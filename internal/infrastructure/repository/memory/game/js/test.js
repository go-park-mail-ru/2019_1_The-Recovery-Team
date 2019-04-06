const { InitEngine } = require("./js");

let cnt = 0;

let initPalyersAction = {
    type: "INIT_PLAYERS",
    payload: JSON.stringify({
        players: [1, 2],
    })
};

let Move = (id, direction) => {
    return {
        type: "INIT_PLAYER_MOVE",
        payload: JSON.stringify({
            player_id: id,
            move: direction,
        })
    }};

const sendInside = InitEngine((actionType, action) => {
    console.log('RECEIVED: ', actionType, 'PAYLOAD: ', JSON.parse(action).payload);
});

console.log(JSON.stringify(Move(20, "DOWN")));

sendInside(JSON.stringify(initPalyersAction));
sendInside(JSON.stringify(Move(1, "DOWN")));
setTimeout(sendInside, 2000, JSON.stringify(Move(1, "RIGHT")));
setTimeout(sendInside, 2000, JSON.stringify(Move(2, "DOWN")));
setTimeout(sendInside, 4000, JSON.stringify(Move(2, "UP")));






