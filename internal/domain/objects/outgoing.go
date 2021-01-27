package objects

import (
	"encoding/binary"
	"math"

	"github.com/steve-nzr/goff/internal/config/constants"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	"github.com/steve-nzr/goff/pkg/abstract"
)

// FPWelcome
const FPWelcomeCmdID = 0x00000000

type FPWelcome struct {
	FPWriter
	ID customtypes.ID
}

func (f *FPWelcome) Serialize() []byte {
	return f.
		Initialize().
		WriteUInt32(FPWelcomeCmdID).
		WriteInt32((int32)(f.ID)).
		finalize()
}

// -- FPWelcome

// FPLoginError
const FPLoginErrorCmdID = 0xfe

type FPLoginError struct {
	FPWriter
	Err error
}

/*
	enum E_LOGIN_ERROR
	{
		ERR_ALREADY_CONNECTED = 0x67, //Account is already connected
		ERR_OLD_CLIENT = 0x6b, //Old client version [closing]
		ERR_CAPACITY = 0x6c, //Capacity has been reached
		ERR_UNAVAILABLE = 0x6d, //Service unavailable
		ERR_BANNED = 0x77, //This ID has been blocked
		ERR_WRONG_PASSWORD = 0x78, //Wrong password
		ERR_WRONG_ID = 0x79, //Wrong ID
		ERR_TIME = 0x80, //Your time is up
		ERR_DATABASE = 0x81, //Other DB error
		ERR_10_PM = 0x83, //You cannot connect after 22:00
		ERR_OUTSIDE_SERVICE = 0x84, //You cannot connect from outside the service for Flyff
		ERR_CHARACTER_CHECKED = 0x85, //Your character is begin checked at the moment
		ERR_WRONG_PASSWORD_AGAIN = 0x86, //Cannot login for 15 seconds
		ERR_WRONG_PASSWORD_WAIT = 0x87, //Cannot login for 15 minutes
		ERR_VERIFICATION = 0x88, //Server verification error
		ERR_SESSION = 0x89, //Session is over

		ERR_NAME_IN_USE = 0x0524
	};
*/
func (f *FPLoginError) getErrorID() int32 {
	switch f.Err {
	case constants.ErrLoginCapacityReached:
		return 0x6c
	case constants.ErrLoginInvalidUserPassword:
		return 0x78
	default:
		return 0x78
	}
}

func (f *FPLoginError) Serialize() []byte {
	return f.
		Initialize().
		WriteUInt32(FPLoginErrorCmdID).
		WriteInt32(f.getErrorID()).
		finalize()
}

// -- FPLoginError

// FPServerList
const FPServerListCmdID = 0xfd

type FPServerList struct {
	FPWriter
	Account string
	AuthKey customtypes.ID
	Servers []*Server
}

func (f *FPServerList) Serialize() []byte {
	p := f.
		Initialize().
		WriteUInt32(FPServerListCmdID).
		WriteUInt32((uint32)(f.AuthKey)).
		WriteByte(1).
		WriteString(f.Account)

	combinedLen := len(f.Servers)
	for i := range f.Servers {
		combinedLen += len(f.Servers[i].Channels)
	}
	p.WriteInt(combinedLen) // wtf

	for i, server := range f.Servers {
		p.WriteUInt32(math.MaxUint32).
			WriteInt32(int32(i + 1)).
			WriteString(server.Name).
			WriteString(server.IP).
			WriteUInt32(0). // b18
			WriteUInt32(0). // lCount
			WriteUInt32(1). // lEnable
			WriteUInt32(0)  // lMax

		for j, channel := range server.Channels {
			p.WriteUInt32(uint32(i + 1)).
				WriteUInt32(uint32(j + 1)).
				WriteString(channel.Name).
				WriteString(channel.IP).
				WriteUInt32(0).                // b18
				WriteUInt32(0).                // lCount
				WriteUInt32(1).                // lEnable
				WriteUInt32(channel.MaxPlayer) // lMax
		}
	}

	return p.finalize()
}

// -- FPServerList

// FPCharacterList
const FPCharacterListCmdID = 0xf3

type FPCharacterList struct {
	FPWriter
	AuthKey    customtypes.ID
	Characters []*entities.Character
}

func (f *FPCharacterList) Serialize() []byte {
	p := f.
		Initialize().
		WriteUInt32(FPCharacterListCmdID).
		WriteUInt32((uint32)(f.AuthKey)).
		WriteInt(len(f.Characters))

	for _, c := range f.Characters {
		p.WriteUInt32((uint32)(c.Slot)).
			WriteUInt32(1).
			WriteUInt32((uint32)(c.Location.MapID)).
			WriteUInt32(0x0B + (uint32)(c.Gender)).
			WriteString(c.Name).
			WriteFloat32((float32)(c.Location.Pos.X)).
			WriteFloat32((float32)(c.Location.Pos.Y)).
			WriteFloat32((float32)(c.Location.Pos.Z)).
			WriteUInt32((uint32)(c.ID)).
			WriteUInt32(0).
			WriteUInt32(0).
			WriteUInt32(0).
			WriteUInt32(c.SkinSetID).
			WriteUInt32(c.HairID).
			WriteUInt32(c.HairColor).
			WriteUInt32(c.HeadID).
			WriteByte(c.Gender).
			WriteUInt32((uint32)(c.JobID)).
			WriteUInt32((uint32)(c.Level)).
			WriteUInt32(0).
			WriteUInt32( /*str*/ 15).
			WriteUInt32( /*sta*/ 15).
			WriteUInt32( /*dex*/ 15).
			WriteUInt32( /*int*/ 15).
			WriteUInt32(0)

		equippedItems := make([]*entities.Item, 0, constants.MaxHumanParts)
		for _, item := range c.Items {
			if item.Position > constants.EquipOffset {
				equippedItems = append(equippedItems, item)
			}
		}
		p.WriteInt(len(equippedItems))
		for _, item := range equippedItems {
			p.WriteInt32(item.ItemID)
		}
	}

	return p.finalize()
}

// -- FPCharacterList

// FPWorldAddress
const FPWorldAddressCmdID = 0xf2

type FPWorldAddress struct {
	FPWriter
	Address string
}

func (f *FPWorldAddress) Serialize() []byte {
	return f.
		Initialize().
		WriteInt32(FPWorldAddressCmdID).
		WriteString(f.Address).
		finalize()
}

// -- FPWorldAddress

// FPPreJoin

const FPPreJoinCmdID = 0xff05

type FPPreJoin struct {
	FPWriter
}

func (f *FPPreJoin) Serialize() []byte {
	return f.
		Initialize().
		WriteInt32(FPPreJoinCmdID).
		finalize()
}

// -- FPPreJoin

// FPMergePacket

type FPMergePacket struct {
	FPWriter
	PacketCount uint16
	MoverID     uint32
}

func (f *FPMergePacket) Initialize(mergeCmd uint32) *FPMergePacket {
	f.initializeEmptyWithSize(32768).
		WriteByte(0x5E). // header
		WriteInt32(0)    // size

	f.PacketCount = 0
	f.
		WriteUInt32(mergeCmd).
		WriteUInt32(0).
		WriteUInt16(f.PacketCount)
	return f
}

func (f *FPMergePacket) Serialize() []byte {
	return f.finalize()
}

func (f *FPMergePacket) finalize() []byte {
	data := f.FPWriter.finalize()
	binary.LittleEndian.PutUint16(data[13:15], f.PacketCount)
	return data
}

func (f *FPMergePacket) AddPacket(moverID customtypes.ID, cmd uint16, p abstract.Serializable) *FPMergePacket {
	f.PacketCount++
	f.
		WriteUInt32((uint32)(moverID)).
		WriteUInt16(cmd).
		WriteBytes(p.Serialize())
	return f
}

// -- FPMergePacket

// FPEnvironmentAll

const FPEnvironmentAllCmdID uint16 = 0x0063

type FPEnvironmentAll struct {
	FPWriter
}

func (f *FPEnvironmentAll) Serialize() []byte {
	return f.initializeEmpty().WriteInt32(0).finalizeWithoutLen()
}

// -- FPEnvironmentAll

// FPAddObj

const FPAddObjCmdID uint16 = 0x00f0

type FPAddObj struct {
	FPWriter
	Character interfaces.Character
}

func (f *FPAddObj) Serialize() []byte {
	f.initializeEmptyWithSize(16384)

	// TODO : remove this wtf twice loop
	for i := 0; i < 2; i++ {
		f.WriteByte(5) // dwObjType && m_dwType -- 5 == OT_MOVER
		// dwObjIndex && m_dwIndex
		if f.Character.GetGender() == 0 {
			f.WriteUInt32(11)
		} else if f.Character.GetGender() == 1 {
			f.WriteUInt32(12)
		}
	}

	// CObj
	f.WriteUInt16(f.Character.GetSize()) // m_vScale
	f.WriteFloat32((float32)(f.Character.GetLocation().X)).
		WriteFloat32((float32)(f.Character.GetLocation().Y)).
		WriteFloat32((float32)(f.Character.GetLocation().Z))

	// CCtrl
	f.WriteInt16(0)                // m_fAngle
	f.WriteID(f.Character.GetID()) // m_objid

	// CMover
	f.WriteUInt16(0)  // m_dwMotion
	f.WriteByte(1)    // m_bPlayer
	f.WriteInt32(230) // m_nHitPoint

	f.WriteInt32(0)  // m_pActMover->AddStateFlag
	f.WriteInt32(0)  // m_pActMover->__ForceSetState
	f.WriteByte(1)   // m_dwBelligerence
	f.WriteInt32(-1) // m_dwMoverSfxId

	f.WriteString(f.Character.GetName())            // m_szName
	f.WriteByte(f.Character.GetGender())            // bySex
	f.WriteByte((byte)(f.Character.GetSkinSetID())) // m_dwSkinSet
	f.WriteByte((byte)(f.Character.GetHairID()))    // m_dwHairMesh
	f.WriteUInt32(f.Character.GetHairColor())       // m_dwHairColor
	f.WriteByte((byte)(f.Character.GetFaceID()))    // m_dwHeadMesh
	f.WriteID(f.Character.GetCharacterID())         // m_idPlayer
	f.WriteByte(f.Character.GetJob())               // m_nJob
	f.WriteUInt16(15)                               // m_nStr
	f.WriteUInt16(15)                               // m_nSta
	f.WriteUInt16(15)                               // m_nDex
	f.WriteUInt16(15)                               // m_nInt
	f.WriteUInt16(f.Character.GetLevel())           // m_nLevel

	f.WriteInt32(0) // m_nFuel
	f.WriteInt32(0) // m_tmAccFuel

	f.WriteByte(0)  // u1 (has guild)
	f.WriteInt32(0) // m_idGuildCloak
	f.WriteByte(0)  // u1 (has party)

	f.WriteByte(100) // m_dwAuthorization
	f.WriteInt32(0)  // m_dwMode
	f.WriteInt32(0)  // m_dwStateMode
	f.WriteInt32(0)  // m_dwUseItemId

	f.WriteInt32(0) // m_dwPKTime
	f.WriteInt32(0) // m_nPKValue
	f.WriteInt32(0) // m_dwPKPropensity
	f.WriteInt32(0) // m_dwPKExp

	f.WriteInt32(0) // m_nFame
	f.WriteByte(0)  // m_nDuel

	f.WriteInt32(-1) // nTemp / m_nHonor

	// m_aEquipInfo[31].nOption
	for i := 0; i < constants.MaxHumanParts; i++ {
		f.WriteInt32(0)
	}

	f.WriteInt32(0) // m_nGuildCombatState

	// m_dwSMTime[26]
	for i := 0; i < 26; i++ {
		f.WriteInt32(0)
	}

	f.WriteInt16(200)                    // m_nManaPoint
	f.WriteInt16(200)                    // m_nFatiguePoint
	f.WriteInt32(0)                      // m_nTutorialState
	f.WriteInt32(0)                      // m_nFxp
	f.WriteUInt32(f.Character.GetGold()) // dwGold

	f.WriteUInt64(0) // m_nExp1
	f.WriteInt32(0)  // m_nSkillLevel
	f.WriteInt32(0)  // m_nSkillPoint
	f.WriteUInt64(0) // m_nDeathExp
	f.WriteInt32(0)  // m_nDeathLevel

	f.WriteInt32(0) // m_idMarkingWorld
	// m_vMarkingPos
	f.WriteFloat32((float32)(f.Character.GetLocation().X)).
		WriteFloat32((float32)(f.Character.GetLocation().Y)).
		WriteFloat32((float32)(f.Character.GetLocation().Z))

	f.WriteByte(0) // m_nQuestSize
	f.WriteByte(0) // m_nCompleteQuestSize
	f.WriteByte(0) // m_nCheckedQuestSize

	f.WriteInt32(0) // m_idMurderer
	f.WriteInt16(0) // m_nRemainGP
	f.WriteInt16(0) // 0

	// m_aEquipInfo[31].dwId
	for i := constants.EquipOffset; i < constants.MaxItems; i++ {
		f.WriteInt32(f.Character.GetItemByPosition((int16)(i)).GetItemID())
	}

	// m_aJobSkill
	for i := 0; i < 45; i++ {
		f.WriteInt32(-1).WriteInt32(0)
	}

	f.WriteByte(0)  // m_nCheerPoint
	f.WriteInt32(0) // m_dwTickCheer

	f.WriteByte(f.Character.GetSlot()) // m_nSlot

	// m_dwGoldBank[constants.MaxCharacters]
	for i := 0; i < constants.MaxCharacters; i++ {
		f.WriteInt32(0)
	}

	// m_idPlayerBank[constants.MaxCharacters]
	for i := 0; i < constants.MaxCharacters; i++ {
		f.WriteInt32(0)
	}

	f.WriteInt32(0) // m_nPlusMaxHitPoint

	f.WriteByte(0) // m_nAttackResistLeft
	f.WriteByte(0) // m_nAttackResistRight
	f.WriteByte(0) // m_nDefenseResist

	f.WriteUInt64(0) // m_nAngelExp
	f.WriteInt32(0)  // m_nAngelLevel

	// inventory
	var size byte
	for i := 0; i < constants.MaxItems; i++ {
		item := f.Character.GetItemByPosition((int16)(i))
		f.WriteInt32((int32)(item.GetUniqueID()))
		if item.GetItemID() != -1 {
			size++
		}
	}
	f.WriteByte(size)
	for i := 0; i < constants.MaxItems; i++ {
		item := f.Character.GetItemByPosition((int16)(i))
		if item.GetItemID() > 0 {
			f.WriteInt8((int8)(item.GetUniqueID())).
				WriteInt32((int32)(item.GetUniqueID()))

			// Item
			f.WriteInt32(item.GetItemID()).
				WriteUInt32(0).
				WriteString("ItemName").
				WriteInt16((int16)(item.GetCount())). // m_nItemNum
				WriteByte(0).                         // m_nRepairNumber
				WriteUInt32(0).                       // m_nHitPoint
				WriteUInt32(0).                       // m_nRepair
				WriteByte(0).                         // m_byFlag
				WriteUInt32(0).                       // m_nAbilityOption
				WriteUInt32(0).                       // m_idGuild
				WriteByte(0).                         // m_bItemResist
				WriteUInt32(0).                       // m_nResistAbilityOption
				WriteUInt32(0).                       // m_nResistSMItemId
				WriteUInt32(0).                       // m_vPiercing.size
				WriteUInt32(0).                       // m_vUltimatePiercing.size
				WriteUInt32(0).                       // SetVisKeepTime.size
				WriteUInt32(0).                       // m_bCharged
				WriteUInt64(0).                       // m_iRandomOptItemId
				WriteUInt32(0).                       // m_dwKeepTime
				WriteByte(0).                         // bPet
				WriteUInt32(0)                        // m_bTranformVisPet
		}
	}
	for i := 0; i < constants.MaxItems; i++ {
		f.WriteInt32((int32)(f.Character.GetItemByPosition((int16)(i)).GetUniqueID()))
	}

	// m_Bank[constants.MaxCharacters]
	for i := 0; i < constants.MaxCharacters; i++ {
		for j := 0; j < constants.BankSize; j++ {
			f.WriteInt(j)
		}
		f.WriteByte(0)
		for j := 0; j < constants.BankSize; j++ {
			f.WriteInt(j)
		}
	}

	f.WriteInt32(-1) // GetPetId()

	// m_Pocket[constants.MaxPockets]
	for i := 0; i < constants.MaxPockets; i++ {
		f.WriteByte(0) // bExists
	}

	f.WriteInt32(0) // m_dwMute

	// m_aHonorTitle[constants.MaxTitle]
	for i := 0; i < constants.MaxTitle; i++ {
		f.WriteInt32(0)
	}

	f.WriteInt32(0) // m_idCampus
	f.WriteInt32(0) // m_nCampusPoint

	f.WriteInt32(0) // m_buffs

	return f.finalizeWithoutLen()
}

// -- FPAddObj

// FPWorldReadInfo

type FPWorldReadInfo struct {
	FPWriter
	Character interfaces.Entity
}

func (f *FPWorldReadInfo) Serialize() []byte {
	f.initializeEmpty()
	return f.
		WriteUInt32((uint32)(f.Character.GetMapID())).
		WriteFloat32((float32)(f.Character.GetLocation().X)).
		WriteFloat32((float32)(f.Character.GetLocation().Y)).
		WriteFloat32((float32)(f.Character.GetLocation().Z)).
		finalizeWithoutLen()
}

// -- FPWorldReadInfo
