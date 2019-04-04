# slugger

My homegrown shell prompt generator


## Add to bash
Write following line to your `.bashrc`:

```bash
export PS1 = '$(slugger)' 
```

## Add to fish
Add function to `~/.config/fish/config.fish`

```fish
function fish_prompt
    slugger
end  
```

## TODOs

- [ ] refactor project
- [ ] add config file support
 
