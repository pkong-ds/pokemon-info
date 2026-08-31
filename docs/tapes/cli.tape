Output docs/cli.gif
Set Shell "zsh"
Set FontSize 15
Set Width 1000
Set Height 480

Hide
Type@30ms "mkdir -p /tmp/pokemon-demo /tmp/pokemon-comp && cp dist/pokemon-info /tmp/pokemon-demo/ && export PATH=/tmp/pokemon-demo:$PATH && pokemon-info completion zsh > /tmp/pokemon-comp/_pokemon-info && fpath=(/tmp/pokemon-comp $fpath) && autoload -U compinit && compinit -u && clear"
Enter
Sleep 800ms
Show

Type@90ms "pokemon-info pi"
Sleep 400ms
Tab
Sleep 1000ms
Enter
Sleep 400ms
Enter
Sleep 3000ms
