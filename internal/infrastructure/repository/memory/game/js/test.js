const { initEngine } = require("./js");

let cnt = 0;

let initPalyersAction = {
    type: "INIT_PLAYERS",
    payload: JSON.stringify({
        playerIds: [1, 2],
    })
};

let Move = (id, direction) => {
    return {
        type: "INIT_PLAYER_MOVE",
        payload: JSON.stringify({
            playerId: id,
            move: direction,
        })
    }
};

let Ready = (id) => {
    return {
        type: "INIT_PLAYER_READY",
        payload: JSON.stringify({
            playerId: id,
        })
    }
};

const sendInside = initEngine((type, payload) => {
    if (payload) {
        console.log('RECEIVED: ', type, 'PAYLOAD: ', payload);
    } else {
        console.log('RECEIVED: ', type);
    }
});

sendInside(JSON.stringify(initPalyersAction));
setTimeout(sendInside, 1000, JSON.stringify(Ready(1)));
setTimeout(sendInside, 1000, JSON.stringify(Ready(2)));
setTimeout(sendInside, 8000, JSON.stringify(Ready(1)));
setTimeout(sendInside, 8000, JSON.stringify(Ready(2)));
// sendInside(JSON.stringify(Move(1, "DOWN")));
// setTimeout(sendInside, 2000, JSON.stringify(Move(1, "RIGHT")));
// setTimeout(sendInside, 2000, JSON.stringify(Move(2, "DOWN")));
// setTimeout(sendInside, 4000, JSON.stringify(Move(2, "UP")));






