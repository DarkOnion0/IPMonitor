{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  shellHook = ''
      # Go command
      go mod verify
      go mod tidy
      go mod download

      # Welcome script
      echo -e "\n$(tput bold)Welcome in the nix-shell for IpChecker$(tput sgr0)"
      
      echo -e "\nList of custom command using 'just' a 'GNU make' like software :"
      echo -e "================================================================"
      just -l
      echo -e "================================================================"
    '';

    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = [
      # Go
      pkgs.go_1_17
      pkgs.golangci-lint
      pkgs.gopls

      # Scripting
      pkgs.zip
      pkgs.unzip
      pkgs.curl
      pkgs.jq
      pkgs.gh

      # Command runner
      pkgs.just

      # Formater
      pkgs.nodePackages.prettier
     ];
}
