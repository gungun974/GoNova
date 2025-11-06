{
  inputs = {
    # Nixpkgs
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    gitignore,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in {
      packages = {
        default = pkgs.buildGo125Module {
          pname = "gonova";
          version = "1.0.0";
          src = gitignore.lib.gitignoreSource ./.;
          subPackages = ["cmd"];
          vendorHash = "sha256-ybpEk5jwAMB4gVd/jjPPYqlzNqCwq0n/9hVIsrzjUFE=";

          flags = [
            "-trimpath"
          ];
          ldflags = [
            "-s"
            "-w"
          ];

          postInstall = ''
            mv $out/bin/cmd $out/bin/gonova
          '';
        };
      };

      devShell = pkgs.mkShell {
        GOROOT = "${pkgs.go_1_25}/share/go";

        buildInputs = with pkgs;
          [
            (golangci-lint.override {buildGo125Module = buildGo125Module;})
            go_1_25
            watchexec
          ]
          ++ lib.optionals stdenv.isDarwin [
            apple-sdk_15
          ];
      };
    });
}
