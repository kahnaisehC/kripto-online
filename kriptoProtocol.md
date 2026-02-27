# how to form messages

- join []
- play [action]
-- action = {"found", "impossible"}
- delete []
- point [player_id]
-- player_id = id of the player being pointed at
- solution [solution]
-- solution = math operation to solve the kripto round at hand
- disconnect 


```
{
    // ID of the user that sends the message
	"Issuer": "number",

    // Type of the message
	// KriptoInvalid = 0
	// KriptoJoin = 1
	// KriptoPlay = 2
	// KriptoDelete = 3 
	// KriptoPoint = 4
	// KriptoSolution = 5
	// KriptoDisconnect = 6
	"Type": "number",

    // if the message was KriptoPlay, sends the type of action to be taken
	// KriptoNil = 0 
	// KriptoFound = 1
	// KriptoImpossible = 2
	"Action": "number",

    // if the message was KriptoPoint, sends the player being pointed at
	"PointedPlayer": "number",

    // if the message was KriptoSolution, sends the player's solution 
	"Solution":     "string" 
}
```

