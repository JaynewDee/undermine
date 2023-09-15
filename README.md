# Undermine
### *A simple terminal-based reminder tool for timing classroom breaks and more*

Current OS support  
`Windows`  
___
#### *Usage*
**Compile Executable (requires a local gotools installation)**  
`cd cmd && go build -o undermine`
*add executable to system path*

__*Run*__
`undermine --note "<reminder_text>" --duration <minutes_til_alarm>`  
___

#### *Features*
- responds to timer end with embedded audio
- responds to timer end with native dialog window
- refreshes time remaining terminal display at minute intervals