{
  description = "X Clone development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs = { nixpkgs.follows = "nixpkgs"; };
    };

    kue.url = "github:Archisman-Mridha/kue";
  };

  outputs = { self, nixpkgs, flake-utils, rust-overlay, kue, }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];

        pkgs = import nixpkgs { inherit system overlays; };
      in with pkgs; {
        devShells.default = mkShell {
          nativeBuildInputs = [ ];

          buildInputs = [
            go_1_24
            golangci-lint

            (python313.withPackages (python313Packages:
              with python313Packages; [
                pip
                venvShellHook
              ]))
            gcc

            llvm
            rustup
            rust-bin.stable.latest.default

            bun
            biome

            buf
            protobuf

            sqlc

            k3d
            cue
            kue.packages.${system}.default
            timoni
            tilt
          ];

          venvDir = "./backend/microservices/models/.venv";

          LD_LIBRARY_PATH = "${stdenv.cc.cc.lib.outPath}/lib:$LD_LIBRARY_PATH";

          shellHook = ''
            if [ ! -d "./backend/microservices/models/.venv" ]; then
              echo "Creating virtual environment"

              python3 -m venv ./backend/microservices/models/.venv
            fi

            source ./backend/microservices/models/.venv/bin/activate
          '';
        };
      });
}
