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
        default = pkgs.buildGo123Module {
          pname = "gonova";
          version = "1.0.0";
          src = gitignore.lib.gitignoreSource ./.;
          subPackages = ["cmd"];
          vendorHash = "sha256-GZ0eq0b0Dui0dLrlNB6O1EF5dX21MxUP+m4L3Rrj6y4=";

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
        GOROOT = "${pkgs.go_1_23}/share/go";

        buildInputs = with pkgs; [
          (golangci-lint.override {buildGo123Module = buildGo123Module;})
          go_1_23
          watchexec
        ];
      };
    });
}
