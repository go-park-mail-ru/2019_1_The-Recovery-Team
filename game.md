# Game documentation

### State

```typescript
interface State {
    Field: Field;
    Players: Array<Player>;
    ActiveItems: {
        [id]: ActiveItem;
    };
    RoundNumber: number;
}
```

### Item

```typescript
interface ActiveItem {
    Type: Itemtype;
    UserId: number;
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
    // Если пустое, то не обновлять
    Info: PlayerInfo | null;
}
```

```typescript
interface PlayerInfo {
    Nickname: string;
    Rating: number;
    Avatar: string;
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

###### INIT_PLAYER_MOVE

```typescript
interface Payload {
    Angle: number;
}
```

#### Item

###### SET_ITEM_START

```typescript
interface Payload {
    Id: number;
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
    ItemType: ItemType;
}
```

###### SET_ITEM_REJECT

```typescript
interface Payload {
    ItemType: ItemType;
}
```

###### INIT_ITEM_USE

```typescript
interface Payload {
    ItemType: ItemType;
}
```



