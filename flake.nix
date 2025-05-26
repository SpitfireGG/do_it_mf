{
  description = "A funny ass Go script to remind of my github contribution chart ( dstreak - daily streak )";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    supportedSystems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];

    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});
  in {
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in rec {
      nanoc = pkgs.buildGoModule {
        pname = "dstreak";
        version = "0.1.0";
        src = ./.;
        vendorHash = "sha256-wyg35Xnw2TJirFCHX6DQY9OaeOBJf+xKnYvXk3AKzDU=";
        buildInputs = [
            


        ];
      };

      default = dstreak;
    });

    devShells = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in {
      default = pkgs.mkShell {
        packages = [
          pkgs.go
          pkgs.gotools # all required tools like linting and formatting are inside the gotools by  default
        pkgs.make



        ];
      };
    });
  };
}
