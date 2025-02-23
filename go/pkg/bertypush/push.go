package bertypush

import (
	"context"

	"go.uber.org/zap"

	"berty.tech/berty/v2/go/internal/messengerdb"
	"berty.tech/berty/v2/go/pkg/errcode"
	"berty.tech/berty/v2/go/pkg/messengertypes"
	"berty.tech/berty/v2/go/pkg/protocoltypes"
)

type EventHandler interface {
	HandleOutOfStoreAppMessage(groupPK []byte, message *protocoltypes.OutOfStoreMessage, payload []byte) (*messengertypes.Interaction, bool, error)
}

type messengerPushReceiver struct {
	logger       *zap.Logger
	pushHandler  PushHandler
	eventHandler EventHandler
	db           *messengerdb.DBWrapper
}

type MessengerPushReceiver interface {
	PushReceive(ctx context.Context, input []byte) (*messengertypes.PushReceive_Reply, error)
}

func NewPushReceiver(pushHandler PushHandler, evtHandler EventHandler, db *messengerdb.DBWrapper, logger *zap.Logger) MessengerPushReceiver {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &messengerPushReceiver{
		logger:       logger,
		pushHandler:  pushHandler,
		eventHandler: evtHandler,
		db:           db,
	}
}

func (m *messengerPushReceiver) PushReceive(ctx context.Context, input []byte) (*messengertypes.PushReceive_Reply, error) {
	clear, err := m.pushHandler.PushReceive(ctx, input)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	i, isNew, err := m.eventHandler.HandleOutOfStoreAppMessage(clear.GroupPublicKey, clear.Message, clear.Cleartext)
	if err != nil {
		return nil, errcode.ErrInternal.Wrap(err)
	}

	accountMuted, conversationMuted, err := m.db.GetMuteStatusForConversation(i.ConversationPublicKey)
	if err != nil {
		accountMuted = true
		conversationMuted = true
	}

	hidePreview := true
	account, err := m.db.GetAccount()
	if err == nil && !account.HidePushPreviews {
		hidePreview = false
	}

	return &messengertypes.PushReceive_Reply{
		Data: &messengertypes.PushReceivedData{
			ProtocolData:      clear,
			Interaction:       i,
			AlreadyReceived:   clear.AlreadyReceived || !isNew,
			ConversationMuted: conversationMuted,
			AccountMuted:      accountMuted,
			HidePreview:       hidePreview,
		},
	}, nil
}
