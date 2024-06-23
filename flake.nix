{
  description = "Tools for reading and manipulating PangYa game files";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        inherit (pkgs) lib stdenv;
        pkgs = nixpkgs.legacyPackages.${system};
        gitCommit = self.dirtyShortRev or self.rev or "";
      in
      {
        packages = rec {
          pangfiles = pkgs.buildGo122Module {
            name = "pangfiles";
            src = self;
            vendorHash = pkgs.lib.fileContents ./go.mod.sri;
            ldflags = [ "-X github.com/pangbox/pangfiles/version.GitCommit=${gitCommit}" ];
            meta = {
              mainProgram = "pang";
            };
          };
          default = pangfiles;
        };
        devShell = pkgs.mkShell {
          packages = [
            pkgs.git
            pkgs.gopls
            pkgs.gotools
            pkgs.go_1_22
            pkgs.gnumake
          ] ++ lib.optional stdenv.isDarwin [ pkgs.macfuse-stubs ];
        };
      }
    );
}
