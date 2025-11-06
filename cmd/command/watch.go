package command

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	"github.com/gungun974/gonova/internal/watch"
	"github.com/spf13/cobra"
)

var (
	watchExtensions      = []string{}
	watchFilters         = []string{}
	watchIgnores         = []string{}
	watchNoDefaultIgnore = false
	processSendInterrupt = false
	processKillDelay     = 0
	processRestartDelay  = 0
	restartDebounceDelay = 0
	templDevMode         bool
)

func init() {
	watchCmd := &cobra.Command{Use: "watcher"}

	watchCmd.AddCommand(watchExecCmd)

	watchExecCmd.Flags().StringSliceVarP(&watchExtensions, "ext", "e", []string{}, "Watch files by extension (with or without the prefix dot)")
	watchExecCmd.Flags().StringSliceVarP(&watchFilters, "filter", "f", []string{}, "Watch files by filename/path using a glob-like pattern")
	watchExecCmd.Flags().StringSliceVarP(&watchIgnores, "ignore", "i", []string{}, "Ignore files by filename/path using a glob-like pattern")

	watchExecCmd.Flags().BoolVar(&watchNoDefaultIgnore, "no-default-ignore", false, "Disable built-in ignores (.git, .DS_Store...)")

	watchExecCmd.Flags().BoolVar(&processSendInterrupt, "send-interrupt", false, "Send interrupt signal before killing the current process")
	watchExecCmd.Flags().IntVar(&processKillDelay, "kill-delay", 500, "Delay in ms after sending interrupt before killing the current process")
	watchExecCmd.Flags().IntVar(&processRestartDelay, "restart-delay", 0, "Delay in ms before restarting the process after stopping")
	watchExecCmd.Flags().IntVar(&restartDebounceDelay, "debounce", 10, "Delay in ms to debounce restarts when files change")

	watchExecCmd.Flags().BoolVar(&templDevMode, "templ-dev-mode", false, "Enable environement variable and run templ watch for extra faster page reload")

	watchExecCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(watchCmd)
}

var watchExecCmd = &cobra.Command{
	Use:   "watchexec",
	Short: "Watch exec a command",
	Run:   WatchExectCmd,
	Args:  cobra.MinimumNArgs(1),
}

func WatchExectCmd(_ *cobra.Command, args []string) {
	if !watchNoDefaultIgnore {
		sep := string(filepath.Separator)

		watchIgnores = append(watchIgnores,
			fmt.Sprintf("**%s.DS_Store", sep),
			"watchexec.*.log",
			"*.py[co]",
			"#*#",
			".#*",
			".*.kate-swp",
			".*.sw?",
			".*.sw?x",
			fmt.Sprintf("**%s.bzr%s**", sep, sep),
			fmt.Sprintf("**%s_darcs%s**", sep, sep),
			fmt.Sprintf("**%s.fossil-settings%s**", sep, sep),
			fmt.Sprintf("**%s.git%s**", sep, sep),
			fmt.Sprintf("**%s.hg%s**", sep, sep),
			fmt.Sprintf("**%s.pijul%s**", sep, sep),
			fmt.Sprintf("**%s.svn%s**", sep, sep),
		)
	}

	var templWatchCmd *exec.Cmd
	var templWatchStderr io.ReadCloser

	if templDevMode {
		os.Setenv("TEMPL_DEV_MODE", "true")
		os.Setenv("TEMPL_DEV_MODE_ROOT", "./tmp/")

		utils.CreateDirectory("tmp")

		templWatchCmd = utils.PrepareCmd("templ",
			[]string{
				"generate",
				"--watch",
				"--log-level", "debug",
				"--open-browser=false",
				"--watch-pattern", `(.+\.templ$)|(.+_templ\.txt$)`,
			},
			".",
		)

		templWatchCmd.Stderr = nil

		var err error

		templWatchStderr, err = templWatchCmd.StderrPipe()
		if err != nil {
			logger.WatcherLogger.Fatalf("Fail to start get templ stdout %v", err)
		}

		err = templWatchCmd.Start()
		if err != nil {
			logger.WatcherLogger.Fatalf("Fail to start templ watch process %v", err)
		}
	}

	var cmd *exec.Cmd

	var ignoreExitMessage bool

	start := func() {
		cmd = utils.PrepareCmd(args[0], args[1:], ".")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c

			if cmd.Process != nil {
				pid := cmd.Process.Pid
				syscall.Kill(-pid, syscall.SIGKILL)
			}

			os.Exit(0)
		}()

		err := cmd.Start()
		if err != nil {
			logger.WatcherLogger.Errorf("Fail to start process %v", err)
			return
		}

		go func() {
			err := cmd.Wait()
			if ignoreExitMessage {
				ignoreExitMessage = false
				return
			}
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					logger.WatcherLogger.Warnf("Process exited unexpectedly with status %d", exitErr.ExitCode())
				} else {
					logger.WatcherLogger.Warnf("Process exited unexpectedly: %v", err)
				}
			} else {
				logger.WatcherLogger.Infof("Process exited with code 0")
			}
		}()
	}

	start()

	restart := func() {
		if cmd != nil && cmd.Process != nil {
			if cmd.ProcessState != nil {
				ignoreExitMessage = !cmd.ProcessState.Exited()
			} else {
				ignoreExitMessage = true
			}

			pid := cmd.Process.Pid
			if processSendInterrupt {
				_ = syscall.Kill(-pid, syscall.SIGINT)

				time.Sleep(time.Duration(processKillDelay) * time.Millisecond)
			}

			if err := cmd.Process.Signal(syscall.Signal(0)); err == nil {
				_ = syscall.Kill(-pid, syscall.SIGKILL)
			}

			_, _ = cmd.Process.Wait()
		} else {
			ignoreExitMessage = false
		}

		time.Sleep(time.Duration(processRestartDelay) * time.Millisecond)

		start()
	}

	var debounceTimer *time.Timer
	var debounceLock sync.Mutex

	debounce := func() {
		debounceLock.Lock()
		defer debounceLock.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		debounceTimer = time.AfterFunc(time.Duration(restartDebounceDelay)*time.Millisecond, func() {
			restart()
		})
	}

	if templDevMode {
		go func() {
			scanner := bufio.NewScanner(templWatchStderr)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.Contains(line, "Generated code") {
					fmt.Fprintln(os.Stderr, line)
				}

				if !strings.Contains(line, "Post-generation event received") {
					continue
				}

				matches := regexp.MustCompile(`needsRestart=(\w+).*needsBrowserReload=(\w+)`).FindStringSubmatch(line)
				if len(matches) == 3 {
					needRestart := matches[1] == "true"
					needBrowserReload := matches[2] == "true"

					if needRestart {
						debounce()
					} else if needBrowserReload {
						resp, err := http.Get("http://127.0.0.1:5174/trigger-refresh")
						if err == nil {
							resp.Body.Close()
						}
					}
				}
			}

			if err := scanner.Err(); err != nil {
				logger.WatcherLogger.Errorf("Failed to scan templ output: %v", err)
			}
		}()
	}

	watch.WatchNewFiles(watch.WatcherConfig{
		Extensions: watchExtensions,
		Filter:     watchFilters,
		Ignore:     watchIgnores,
	}, func() {
		debounce()
	})
}
