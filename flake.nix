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
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = rec {
          pangfiles = pkgs.buildGoModule {
            pname = "pangfiles";
            version = "0.0.0";
            src = self;
            vendorHash = "sha256-LwSUYQ+mt2dlyEMwlwCI/OZR8EM1jXLKfbi0w0zWgDM=";
            meta = {
              mainProgram = "pang";
            };
          };
          default = pangfiles;
        };
      }
    );
}
