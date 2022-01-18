(function () {
  window.state = {
    Board: {
      Rows: 4,
      Cols: 6,
      PieceCoords: {
        TL: [
          { Row: 0, Col: 0 },
          { Row: 0, Col: 1 },
          { Row: 1, Col: 0 },
        ],
        TR: [
          { Row: 0, Col: 2 },
          { Row: 0, Col: 3 },
          { Row: 1, Col: 3 },
        ],
        BL: [
          { Row: 2, Col: 0 },
          { Row: 3, Col: 0 },
          { Row: 3, Col: 1 },
        ],
        BR: [
          { Row: 2, Col: 3 },
          { Row: 3, Col: 2 },
          { Row: 3, Col: 3 },
        ],
        SQ: [
          { Row: 1, Col: 4 },
          { Row: 1, Col: 5 },
          { Row: 2, Col: 4 },
          { Row: 2, Col: 5 },
        ],
      },
    },
    PrecedingMoves: [
      { Piece: "SQ", Dir: "up" },
      { Piece: "BR", Dir: "right" },
      { Piece: "BL", Dir: "right" },
      { Piece: "BR", Dir: "right" },
      { Piece: "TL", Dir: "down" },
      { Piece: "TR", Dir: "down" },
      { Piece: "TR", Dir: "down" },
      { Piece: "SQ", Dir: "left" },
      { Piece: "BR", Dir: "up" },
      { Piece: "SQ", Dir: "left" },
      { Piece: "BR", Dir: "up" },
      { Piece: "TR", Dir: "right" },
      { Piece: "TR", Dir: "right" },
      { Piece: "BL", Dir: "right" },
      { Piece: "BL", Dir: "right" },
      { Piece: "TL", Dir: "down" },
      { Piece: "SQ", Dir: "left" },
      { Piece: "BR", Dir: "left" },
      { Piece: "SQ", Dir: "left" },
      { Piece: "BR", Dir: "left" },
      { Piece: "TR", Dir: "up" },
      { Piece: "TR", Dir: "up" },
      { Piece: "BL", Dir: "right" },
      { Piece: "BL", Dir: "up" },
      { Piece: "BR", Dir: "down" },
      { Piece: "BR", Dir: "down" },
      { Piece: "BR", Dir: "left" },
      { Piece: "BL", Dir: "left" },
      { Piece: "BL", Dir: "up" },
      { Piece: "BR", Dir: "right" },
      { Piece: "BR", Dir: "right" },
      { Piece: "BR", Dir: "right" },
      { Piece: "BL", Dir: "down" },
      { Piece: "BL", Dir: "left" },
      { Piece: "BL", Dir: "down" },
      { Piece: "SQ", Dir: "right" },
      { Piece: "SQ", Dir: "right" },
      { Piece: "TL", Dir: "up" },
      { Piece: "TL", Dir: "up" },
      { Piece: "BL", Dir: "left" },
      { Piece: "SQ", Dir: "down" },
      { Piece: "TL", Dir: "right" },
      { Piece: "TR", Dir: "left" },
      { Piece: "BR", Dir: "left" },
    ],
  };

  window.COLORS = {
    TR: randomRGB(),
    TL: randomRGB(),
    BR: randomRGB(),
    BL: randomRGB(),
    SQ: randomRGB(),
  };

  window.allMoves = state.PrecedingMoves;
  window.STEP_IDX = 0;

  window.stepBack = function stepBack() {
    if (STEP_IDX <= 0) return alert("Already at initial state");

    let step = allMoves[--STEP_IDX];
    applyStep(step);
    draw();
  };

  window.stepFwd = function stepFwd() {
    if (STEP_IDX > allMoves.length - 1) return alert("Already at final state");

    applyStep(allMoves[STEP_IDX++]);
    draw();
  };

  function applyStep({ Piece, Dir }) {
    let coords = state.Board.PieceCoords;
    switch (Dir) {
      case "up":
        coords[Piece] = coords[Piece].map(({ Row, Col }) => ({
          Row: Row - 1,
          Col: Col,
        }));
        break;
      case "down":
        coords[Piece] = coords[Piece].map(({ Row, Col }) => ({
          Row: Row + 1,
          Col: Col,
        }));
        break;
      case "left":
        coords[Piece] = coords[Piece].map(({ Row, Col }) => ({
          Row: Row,
          Col: Col - 1,
        }));
        break;
      case "right":
        coords[Piece] = coords[Piece].map(({ Row, Col }) => ({
          Row: Row,
          Col: Col + 1,
        }));
        break;
      default:
        alert("unexpected direction! " + Dir);
    }
  }
  window.draw = function draw() {
    [...document.getElementsByClassName("row")].forEach((r) =>
      [...r.children].forEach(
        (child) => (child.style.backgroundColor = "white")
      )
    );
    let pieceCoords = state.Board.PieceCoords;
    let pieces = Object.keys(pieceCoords);
    pieces.forEach((piece) => {
      pieceCoords[piece].forEach(({ Row, Col }) => {
        let row = document.getElementsByClassName(`row-${Row}`)[0];
        let cell = row.children[Col];
        if (!cell) debugger;
        cell.style.backgroundColor = rgbToString(COLORS[piece]);
      });
    });
  };

  function randomRGB() {
    return [
      Math.round(Math.random() * 0xff),
      Math.round(Math.random() * 0xff),
      Math.round(Math.random() * 0xff),
    ];
  }
  function rgbToString([r, g, b]) {
    return `rgb(${r}, ${g}, ${b})`;
  }

  draw();
})();
