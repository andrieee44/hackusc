{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      system = "x86_64-linux";

      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          (final: prev: {
            go-migrate = prev.go-migrate.overrideAttrs (oldAttrs: {
              tags = [ "mysql" ];
            });
          })
        ];
      };
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "hackusc";
        version = "0.0.1";
        vendorHash = null;
        src = self;
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          go-migrate
          lsof
          mariadb
          openssh
          sqlc
        ];
      };
    };
}
