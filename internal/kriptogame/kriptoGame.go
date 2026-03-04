package kriptogame

import (
	"errors"
	"strconv"
	"strings"
)

type Phase int

const (
	PhasePending Phase = iota
	PhaseWaitingActions
	PhaseWaitingPointer
	PhaseWaitingSolution
	PhaseFinished
)

type PlayerState int

const (
	PlayerStatePending PlayerState = iota
	PlayerStateDefeated
	PlayerStateImpossible
	PlayerStateFound
	PlayerStatePoint
	PlayerStatePointed
)

type Game struct {
	// Static Information

	// Dynamic Information
	PlayersState []PlayerState
	PlayerOrder  []int

	// Game Information
	Turn       int
	Phase      Phase
	PointedIdx int
	Cards      []Card
}

func NewGame(n int) Game {
	game := Game{
		PlayersState: make([]PlayerState, n),
		PlayerOrder:  make([]int, n),
		Turn:         0,
		Phase:        PhaseWaitingActions,
		PointedIdx:   -1,
		Cards:        generateHand(),
	}

	return game
}

var errorBadFormat = errors.New("bad format")

type MessageType int

const (
	Invalid MessageType = iota
	TypeStart
	TypeJoin
	TypePlay
	TypeDelete
	TypePoint
	TypeSolution
	TypeDisconnect
)

type Action int

const (
	ActionNil = iota
	ActionFound
	ActionImpossible
)

type KriptoMessage struct {
	IssuerIdx     int
	Type          MessageType
	Action        Action
	PointedPlayer int
	Solution      string
}

func (game *Game) CheckSolution(solution string) bool {
	expressions := strings.Split(solution, ",")
	if len(expressions) != 3 {
		return false
	}

	values := make([]int, 0, 4)
	for _, c := range game.Cards {
		values = append(values, c.Value)
	}

	for _, exp := range expressions {
		splittedExp := strings.Split(exp, ";")
		if len(splittedExp) != 3 {
			return false
		}
		val1, err := strconv.Atoi(splittedExp[1])
		if err != nil {
			return false
		}
		val2, err := strconv.Atoi(splittedExp[2])
		if err != nil {
			return false
		}
		var res int
		switch splittedExp[0] {
		case "+":
			res = val1 + val2
		case "*":
			res = val1 * val2
		case "-":
			res = val1 - val2
			if res < 0 {
				res = -res
			}
		case "/":
			if val2 == 0 {
				return false
			}
			if val1%val2 != 0 {
				return false
			}
			res = val1 / val2
		}
		newVals := make([]int, 0, len(values)-1)
		for _, v := range values {
			if v == val1 {
				val1 = -1
				continue
			}
			if v == val2 {
				val2 = -1
				continue
			}
			newVals = append(newVals, v)
		}
		if val1 != -1 || val2 != -1 {
			return false
		}
		newVals = append(newVals, res)
		values = newVals
	}
	if len(values) != 1 {
		return false
	}
	if values[0] != game.Cards[len(game.Cards)-1].Value {
		return false
	}

	return true
}

func (game *Game) ParseMessage(msg string) (KriptoMessage, error) {
	kriptoMsg := KriptoMessage{}
	kriptoMsg.PointedPlayer = -1

	splittedMsg := strings.Split(msg, " ")
	if len(splittedMsg) <= 1 {
		return kriptoMsg, errorBadFormat
	}

	issuerIdxString := splittedMsg[0]
	issuerIdx, err := strconv.Atoi(issuerIdxString)
	if err != nil {
		return kriptoMsg, err
	}
	// TODO: Erase this?
	if issuerIdx < 0 || issuerIdx >= len(game.PlayersState) {
		return kriptoMsg, errors.New("issuer Idx out of range")
	}

	msgType := splittedMsg[1]

	switch msgType {

	case "play":
		if len(splittedMsg) < 3 {
			return kriptoMsg, errorBadFormat
		}
		action := splittedMsg[2]
		switch action {
		case "found":
			kriptoMsg.Action = ActionFound
		case "impossible":
			kriptoMsg.Action = ActionImpossible
		default:
			return kriptoMsg, errorBadFormat
		}
		kriptoMsg.Type = TypePlay

	case "point":
		if len(splittedMsg) < 3 {
			return kriptoMsg, errorBadFormat
		}
		pointedIdx, err := strconv.Atoi(splittedMsg[2])
		if err != nil {
			return kriptoMsg, err
		}
		// TODO: Erase this?
		if pointedIdx < 0 || pointedIdx >= len(game.PlayersState) {
			return kriptoMsg, errors.New("pointedIdx out of range")
		}
		if game.PlayersState[pointedIdx] != PlayerStateFound {
			return kriptoMsg, errors.New("pointed player is not able to be pointed")
		}
		kriptoMsg.PointedPlayer = pointedIdx

	case "solution":
		/// AAAAAAA
		if len(splittedMsg) < 3 {
			return kriptoMsg, errorBadFormat
		}
		if ok := game.CheckSolution(splittedMsg[2]); !ok {
			return kriptoMsg, errors.New("invalid solution")
		}
		kriptoMsg.Solution = splittedMsg[2]
		kriptoMsg.Type = TypeSolution

	case "disconnect":
		kriptoMsg.Type = TypeDisconnect

	default:
		return kriptoMsg, errors.New("message Type not supported" + msgType)
	}

	kriptoMsg.IssuerIdx = issuerIdx
	return kriptoMsg, nil
}

func (game *Game) CheckMessageValidity(msg KriptoMessage) error {
	if msg.IssuerIdx < 0 || msg.IssuerIdx >= len(game.PlayersState) {
		return errors.New("issuerIdx out of range")
	}
	switch msg.Type {
	case TypePlay:
		if game.PlayersState[msg.IssuerIdx] != PlayerStatePending {
			return errors.New("player is not pending")
		}
		if game.Phase != PhasePending {
			return errors.New("game is not in Pending state")
		}
		switch msg.Action {
		case ActionFound:
		case ActionImpossible:
		default:
			return errors.New("action is not supported")
		}

	case TypePoint:
		if game.Phase != PhaseWaitingPointer {
			return errors.New("not waiting pointer")
		}
		if game.PlayersState[msg.IssuerIdx] != PlayerStatePoint {
			return errors.New("issuer is not the pointer player")
		}
		if msg.PointedPlayer < 0 || msg.PointedPlayer >= len(game.PlayersState) {
			return errors.New("pointed player idx out of range")
		}
		if game.PlayersState[msg.PointedPlayer] != PlayerStateFound {
			return errors.New("cannot point to that player")
		}

	case TypeSolution:
		if game.PlayersState[msg.IssuerIdx] != PlayerStatePointed {
			return errors.New("the issuer is not being pointed")
		}
		if game.Phase != PhaseWaitingSolution {
			return errors.New("the game is not waiting for a solution")
		}

	case TypeDisconnect:
		if game.Phase == PhaseFinished {
			return errors.New("game already ended")
		}
	default:
		// case TypeDelete:
		// case TypeStart: //
		// case TypeJoin: //
		return errors.New("unsupported type")
	}
	return nil
}

func (game *Game) ExecuteUnsafe(msg KriptoMessage) bool {
	// assumes the message is already ok (aka VerifyMessageValidity returned nil

	game.Turn++

	switch msg.Type {
	case TypePlay:
		switch msg.Action {
		case ActionFound:
			game.PlayersState[msg.IssuerIdx] = PlayerStateFound
		case ActionImpossible:
			game.PlayersState[msg.IssuerIdx] = PlayerStateImpossible
		}
		var cntPending int
		pointer := -1
		for i, s := range game.PlayersState {
			if s == PlayerStatePending {
				cntPending++
				pointer = i
			}
		}
		if cntPending == 1 {
			game.PlayersState[pointer] = PlayerStatePoint
			game.Phase = PhaseWaitingPointer
		}

	case TypePoint:
		game.PlayersState[msg.PointedPlayer] = PlayerStatePointed
		game.Phase = PhaseWaitingSolution

	case TypeSolution:
		if game.CheckSolution(msg.Solution) {
			for i, s := range game.PlayersState {
				if s == PlayerStateImpossible || s == PlayerStatePending || s == PlayerStatePoint {
					game.PlayersState[i] = PlayerStateDefeated
				}
			}
		} else {
			game.PlayersState[msg.IssuerIdx] = PlayerStateDefeated
		}
		game.Cards = generateHand()
		game.Phase = PhaseWaitingActions

	case TypeDisconnect:
		game.PlayersState[msg.IssuerIdx] = PlayerStateDefeated
	default:
		// case TypeStart: //
		// case TypeJoin: //
		// case TypeDelete:
		panic("ExecuteUnsafe executed a bad command")
	}

	cntAlivePlayers := 0
	winner := -1
	for i, s := range game.PlayersState {
		if s != PlayerStateDefeated {
			cntAlivePlayers++
			winner = i
		}
	}
	_ = winner
	if cntAlivePlayers <= 1 {
		game.Phase = PhaseFinished
	}

	return true
}
