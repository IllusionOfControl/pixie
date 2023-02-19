package bot

type State int

const (
	StartState State = iota

	PixilizerAskedForCount
	PixilizerAskedForImage

	PalettizerAskedForDimension
	PalettizerAskedForImage
)

type UserState struct {
	State   State
	Context map[string]string
}

func (u UserState) ClearContext() {
	u.Context = make(map[string]string)
}

var userStates = make(map[int64]*UserState)

func GetUserState(userID int64) *UserState {
	userState, ok := userStates[userID]
	if !ok {
		userState = &UserState{State: StartState, Context: make(map[string]string)}
		userStates[userID] = userState
	}
	return userState
}
