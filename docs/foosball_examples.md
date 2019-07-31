# Foosball examples
```
1. Regular goal
    POST {"scorer": "user1", "opponent": "user2", "player": "p3", "gamelle": false}
    Returns:
        {"user1": {"sets": 0, "points": 1}, "user2": {"sets": 0, "points": 0}, "goals_in_balance": 0}
2. A gamelle is scored
    POST {"scorer": "user1", "opponent": "user2", "player": "p3", "gamelle": true}
    Returns:
        {"user1": {"sets": 0, "points": 0}, "user2": {"sets": 0, "points": -1}, "goals_in_balance": 0}
3. Demi + goal
    POST {"scorer": "user1", "opponent": "user2", "player": "p6", "gamelle": false}
    Returns:
        {"user1": {"sets": 0, "points": 0}, "user2": {"sets": 0, "points": 0}, "goals_in_balance": 2}
    POST {"scorer": "user2", "opponent": "user1", "player": "p3", "gamelle": false}
    Returns:
        {"user1": {"sets": 0, "points": 0}, "user2": {"sets": 0, "points": 2}, "goals_in_balance": 0}
```