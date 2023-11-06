{ pkgs, ... }:

{
  # https://devenv.sh/basics/
  # env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = with pkgs; [ git helix gcc ];

  # https://devenv.sh/scripts/
  # scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hx -g fetch && hx -g build
  '';

  # https://devenv.sh/languages/
  # languages.nix.enable = true;
  languages.go = {
    enable = true;
    package = pkgs.go_1_21;
  };

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
 }
