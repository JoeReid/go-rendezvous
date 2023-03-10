package rendezvous

import (
	"testing"

	mockhash "github.com/JoeReid/go-rendezvous/mock_hash"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHasher_Members(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	assert.Empty(t, hasher.Members())

	hasher.AddMembers("node1", "node2")
	assert.Equal(t, []string{"node1", "node2"}, hasher.Members())

	hasher.RemoveMembers("node1")
	assert.Equal(t, []string{"node2"}, hasher.Members())

	hasher.RemoveMembers("node2")
	assert.Empty(t, hasher.Members())

	hasher.SetMembers("node1", "node2")
	assert.Equal(t, []string{"node1", "node2"}, hasher.Members())
}

func TestHasher_Prioritise_no_members(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	assert.Empty(t, hasher.Prioritise("item"))
}

func TestHasher_Prioritise_one_member(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	hasher.AddMembers("node1")

	assert.Equal(t, []string{"node1"}, hasher.Prioritise("item"))
}

func TestHasher_Place_no_members(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	assert.Empty(t, hasher.Place("item", 0))
	assert.Empty(t, hasher.Place("item", 1))
	assert.Empty(t, hasher.Place("item", 2))
}

func TestHasher_Place_one_member(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	hasher.AddMembers("node1")

	assert.Empty(t, hasher.Place("item", 0))
	assert.Equal(t, []string{"node1"}, hasher.Place("item", 1))
	assert.Equal(t, []string{"node1"}, hasher.Place("item", 2))
	assert.Equal(t, []string{"node1"}, hasher.Place("item", 3))
}

func TestHasher_Prioritise_two_members(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	gomock.InOrder(
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x01}),
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x02}),
	)

	hasher.AddMembers("node1", "node2")

	assert.Equal(t, []string{"node1", "node2"}, hasher.Prioritise("item"))
}

func TestHasher_Place_two_members_0(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	hasher.AddMembers("node1", "node2")

	assert.Empty(t, hasher.Place("item", 0))
}

func TestHasher_Place_two_members_1(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	gomock.InOrder(
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x01}),
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x02}),
	)

	hasher.AddMembers("node1", "node2")

	assert.Equal(t, []string{"node1"}, hasher.Place("item", 1))
}

func TestHasher_Place_two_members_2(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	gomock.InOrder(
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x01}),
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x02}),
	)

	hasher.AddMembers("node1", "node2")

	assert.Equal(t, []string{"node1", "node2"}, hasher.Place("item", 2))
}

func TestHasher_Place_two_members_3(t *testing.T) {
	var (
		m      = mockhash.NewMockHash(gomock.NewController(t))
		hasher = NewHasher(WithHashImplementation(m))
	)

	gomock.InOrder(
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x01}),
		m.EXPECT().Reset(),
		m.EXPECT().Write(gomock.Any()).Return(5, nil).Times(2),
		m.EXPECT().Sum(nil).Return([]byte{0x00, 0x00, 0x00, 0x02}),
	)

	hasher.AddMembers("node1", "node2")

	assert.Equal(t, []string{"node1", "node2"}, hasher.Place("item", 3))
}
