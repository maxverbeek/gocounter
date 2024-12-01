{
  inputs = {
    nixpkgs.url = "nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShell = pkgs.mkShell {
          name = "devshell";
          packages = with pkgs; [
            go_1_23
            golangci-lint
            cobra-cli
          ];
        };

        overlays.default = final: prev: {
          gocount = final.buildGoModule {
            pname = "gocount";
            version = "0.1.0";
            src = ./.;
            subPackages = [ "cmd/gocount" ];
            vendorHash = null;
          };
        };

        packages.default =
          (import nixpkgs {
            inherit system;
            overlays = [ self.overlays.${system}.default ];
          }).gocount;
      }
    );
}
