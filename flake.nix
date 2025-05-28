{
  description = "Is An AI [Bifrost Client] Environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }: 
  let
    pkgs = import nixpkgs { system = "aarch64-darwin"; };
    go = pkgs.go;
    npm = pkgs.nodejs_24;
    wails = pkgs.wails;
  in {
    devShells.aarch64-darwin.default = pkgs.mkShell {
      buildInputs = [
        go
        npm
        wails
      ];
      
      shellHook = ''
        echo "Is An AI [Bifrost Client] Environment"
        echo "Enter Go $(${go}/bin/go --version)"
        echo "Enter NPM $(${npm}/bin/npm --version)"
        echo "Enter Wails"
      '';
    };
  };
}

