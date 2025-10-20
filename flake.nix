{
  description = "X Clone development environment";

  inputs = {
    nixpkgs.url = "github:cachix/devenv-nixpkgs/rolling";
    devenv.url = "github:cachix/devenv";

    flake-utils.url = "github:numtide/flake-utils";

    nixpkgs-python = {
      url = "github:cachix/nixpkgs-python";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    kue.url = "github:Archisman-Mridha/kue";
  };

  nixConfig = {
    extra-substituters = "https://devenv.cachix.org";

    extra-trusted-public-keys =
      "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
  };

  outputs = { self, nixpkgs, kue, ... }@inputs:
    inputs.flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import inputs.rust-overlay) ];

        pkgs = import nixpkgs { inherit system overlays; };
      in with pkgs; {
        devShells.default = inputs.devenv.lib.mkShell {
          inherit inputs pkgs;

          modules = [
            ({ pkgs, config, ... }: {
              # Most packages come pre-built with binaries provided by the official Nix binary
              # cache. If you're modifying a package or using a package that's not built
              # upstream, Nix will build it from source instead of downloading a binary.
              # To prevent packages from being built more than once, devenv provides seamless
              # integration with binary caches hosted by Cachix.
              #
              # Devenv will automatically configure Cachix caches for you, or guide you how to
              # add the caches to Nix manually. Any caches set up by devenv are used in addition
              # to the caches configured in Nix, for example, in /etc/nix/nix.conf.
              cachix = {
                enable = true;

                # devenv.cachix.org is added to the list of pull caches by default. It mirrors
                # the official NixOS cache and is designed to provide caching for the
                # devenv-nixpkgs/rolling nixpkgs input.
                #
                # Some languages and integrations may automatically add caches when enabled.
                pull = [ ];
              };

              languages = {
                java = {
                  enable = true;

                  gradle.enable = true;
                };

                python = {
                  enable = true;
                  version = "3.11";

                  venv = {
                    enable = true;
                    requirements =
                      ./backend/microservices/models/requirements.txt;
                  };
                };
              };

              packages = with pkgs; [
                go
                golangci-lint

                buf
                protobuf

                sqlc

                gcc

                llvm
                rustup
                rust-bin.stable.latest.default

                bun
                biome

                cue
                k3d
                kue.packages.${system}.default
                timoni
                tilt
              ];
            })
          ];
        };
      });
}
