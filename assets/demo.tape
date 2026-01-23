# itodo demo — requires ./itodo (run: make build)
Output demo.gif
Set FontSize 13
Set Width 880
Set Height 520
Set Padding 20
Set Theme "Monokai Pro"

# Isolated home for a clean, reproducible demo
Type "./itodo"
Enter

Sleep 1s
# Create Todo 1
Type "n"
Sleep 500ms
Type "Todo 1"
Tab
Sleep 200ms
Ctrl+S
Sleep 600ms
# Create Todo 2
Type "n"
Sleep 400ms
Type "Todo 2"
Tab
Sleep 150ms
Ctrl+S
Sleep 500ms
# Create Todo 3
Type "n"
Sleep 400ms
Type "Todo 3"
Tab
Sleep 150ms
Ctrl+S
Sleep 500ms
# Create Todo 4
Type "n"
Sleep 400ms
Type "Todo 4"
Tab
Sleep 150ms
Ctrl+S
Sleep 600ms
# Move cursor to Todo 3 (up from Todo 4)
Up 1
Sleep 300ms
# Indent Todo 3 as subtask of Todo 1
Type ">"
Sleep 500ms
# Move cursor to Todo 4 (down 2: past Todo 3 child and Todo 2)
Down 2
Sleep 300ms
# Indent Todo 4 as subtask of Todo 1
Type ">"
Sleep 600ms
# Toggle Todo 1 complete to show the feature
Up 3
Sleep 300ms
Enter
Sleep 400ms
# Switch to next date
Right
Sleep 500ms
# Quick switch back to today
Space
Sleep 500ms
# Show calendar
Type ";"
Sleep 800ms
# Move to another date in calendar (right arrow moves day, down moves week)
Right 3
Sleep 400ms
Down
Sleep 400ms
# Select the date
Enter
Sleep 500ms
# Quit
Type "q"
