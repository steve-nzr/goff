package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/pkg/testutils/mock_interfaces"
	"github.com/stretchr/testify/assert"
)

func TestWelcomeUseCase_Greet(t *testing.T) {
	// mock init
	ctrl := gomock.NewController(t)
	mockIdGen := mock_interfaces.NewMockIdentifierGenerator(ctrl)

	uc := NewWelcome(mockIdGen)

	// expect
	mockIdGen.EXPECT().
		Generate().
		Return((customtypes.ID)(5))

	// then
	id, res := uc.Greet()
	assert.EqualValues(t, 5, id)
	assert.Equal(t, &objects.FPWelcome{
		ID: 5,
	}, res)
}
