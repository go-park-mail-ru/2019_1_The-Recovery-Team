# Game documentation

### State

```typescript
interface State {
    Field: Field;
    Players: { 
        [id: string]: Player;
    };
    ActiveItems: {
        [id]: ActiveItem;
    };
    RoundNumber: number;
    RoundTimer: number;
}
```

### Item

```typescript
interface ActiveItem {
    Type: Itemtype;
    PlayerId: number;
    Duration: number;
}
```

### Field

```typescript
interface Field {
    Cells: Array<Cell>;
    Width: number;
    Height: number;
}
```

### Cell

```typescript
interface Cell {
    Row: number;
    Col: number;
    Type: CellType;
    HasBox: boolean;
}
```

###### CellType

```typescript
enum CellType {
    WATER = 'WATER',
    SAND = 'SAND',
    SWAMP = 'SWAMP',
}
```

### Player(Крабик)

```typescript
interface Player {
    Id: int;
    X: int;
    Y: int;
    Items: {
        [Type: ItemType]: number;
    }
    LoseRound: number | null;
}
```

### ItemType

```typescript
enum ItemType {
    LIFEBUOY = 'LIFEBUOY',
    BOMB = 'BOMB',
    SAND = 'SAND'
}
```
# Actions

```typescript
interface Action {
    Type: ActionType;
    PayLoad: {
        id,
        data
    };
}
```
### Action types

#### Engine

###### INIT_ENGINE_STOP

```typescript
interface Payload {}
```

###### SET_ENGINE_STOP

```typescript
interface Payload {}
```

#### Game

###### SET_GAME_START

```typescript
interface Payload {
    Field: Field;
    Players: Array<Player>;
}
```

###### SET_GAME_STOP

```typescript
interface Payload {
    Players: Array<Player>;
}
```

###### INIT_GAME_PAUSE

```typescript
interface Payload {}
```

###### SET_GAME_PAUSE

```typescript
interface Payload {}
```

###### INIT_GAME_RESUME
```typescript
interface Payload {}
```

###### SET_GAME_RESUME
```typescript
interface Payload {}
```
###### INIT_GAME_QUIT
```typescript
interface Payload {}
```

#### Round

###### SET_ROUND_START

```typescript
interface Payload {
    RoundNumber: number;
    Duration: number;
}
```

###### SET_ROUND_STOP

```typescript
interface Payload {}
```

#### Field

###### SET_FIELD_DIFF

```typescript
interface Payload {
    Cells: Arrray<Cell>;
}
```

#### Player

###### SET_PLAYER

```typescript
interface Payload {
    Player: Player;
}
```

###### INIT_PLAYERS
Frontend only!
```typescript
interface Payload {
    PlayerIds: Array<number>;
}
```

###### INIT_PLAYER_READY
```typescript
interface Payload {
    PlayerId: number;
}
```

###### INIT_PLAYER_MOVE

```typescript
interface Payload {
    PlayerId: number;
    Move: string;
}
```

#### Item

###### SET_ITEM_START

```typescript
interface Payload {
    Id: string;
    Item: ActiveItem;
}
```

###### SET_ITEM_STOP

```typescript
interface Payload {
    Id: number;
}
```

###### SET_ITEM_REJECT

```typescript
interface Payload {
    PlayerId: number;
    ItemType: ItemType;
}
```

###### INIT_ITEM_USE

```typescript
interface Payload {
    PlayerId: number;
    ItemType: ItemType;
}
```
