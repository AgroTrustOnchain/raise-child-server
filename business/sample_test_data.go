package business

import (
	"raise-child/constants/shared"
	"raise-child/model/entities"
)

// Valid address
const (
	sampleAddress        string = "0x1088c00ff68478ee818c97687be5b485c4274211795763b3f021ed2c05e8fc4d"
	sampleInvalidAddress string = "abc"
	sampleSub            string = "12345678901234567890"
)

// Manage Object
var (
	sampleManageObj = entities.Manage{
		ID:                    entities.ID{ID: sampleManageObjId},
		AdminIds:              sampleAdminIds,
		AdminNfts:             sampleAdminNfts,
		ChildIds:              sampleChildIds,
		VolunteerIds:          sampleVolunteerIds,
		VolunteerNfts:         sampleVolunteerNfts,
		LocalLeaderIds:        sampleLocalLeaderIds,
		LocalLeaderNfts:       sampleLocalLeaderNfts,
		LocalRegions:          sampleLocalRegions,
		ChildrenCenters:       sampleChildrenCenters,
		CenterConfirmStatuses: sampleCenterConfirmStatuses,
		CreatedCenters:        sampleCreatedCenters,
		DonorIds:              sampleDonorIds,
		DonorNfts:             sampleDonorNfts,
	}
	sampleJsonManageObj = map[string]interface{}{
		"id": map[string]string{
			"id": sampleManageObjId,
		},
		"admin_ids":               sampleAdminIds,
		"admin_nfts":              sampleAdminNfts,
		"child_ids":               sampleChildIds,
		"volunteer_ids":           sampleVolunteerIds,
		"volunteer_nfts":          sampleVolunteerNfts,
		"local_leader_nfts":       sampleLocalLeaderNfts,
		"local_leader_ids":        sampleLocalLeaderIds,
		"local_regions":           sampleLocalRegions,
		"children_centers":        sampleChildrenCenters,
		"center_confirm_statuses": sampleCenterConfirmStatuses,
		"created_centers":         sampleCreatedCenters,
		"donor_ids":               sampleDonorIds,
		"donor_nfts":              sampleDonorNfts,
	}
)

// Admin NFT
var (
	sampleJsonAdminNft1 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[0],
		},
		"owner":                 sampleAdminIds[0],
		"identity_code":         "001092000123",
		"identity_card_blob_id": "blob_v1_88291",
		"avatar_blob_id":        "avatar_v1_112",
		"first_name":            "Minh",
		"last_name":             "Nguyễn Văn",
		"gender":                "male",
		"date_of_birth":         "15/05/1992",
		"phone_number":          "094901234567",
		"email":                 "minh.nguyen@example.com",
		"uploaded_at":           "1698397200000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft2 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[1],
		},
		"owner":                 sampleAdminIds[1],
		"identity_code":         "042095000456",
		"identity_card_blob_id": "blob_v1_99302",
		"avatar_blob_id":        "avatar_v1_223",
		"first_name":            "Linh",
		"last_name":             "Trần Thị",
		"gender":                "Female",
		"date_of_birth":         "1995-08-20",
		"phone_number":          "+84912345678",
		"email":                 "linh.tran@example.com",
		"uploaded_at":           "1698397500000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft3 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[2],
		},
		"owner":                 sampleAdminIds[2],
		"identity_code":         "079088000789",
		"identity_card_blob_id": "blob_v1_44567",
		"avatar_blob_id":        "avatar_v1_334",
		"first_name":            "Hoàng",
		"last_name":             "Lê Anh",
		"gender":                "Male",
		"date_of_birth":         "1988-12-12",
		"phone_number":          "+84923456789",
		"email":                 "hoang.le@example.com",
		"uploaded_at":           "1698397800000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft4 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[3],
		},
		"owner":                 sampleAdminIds[3],
		"identity_code":         "001099000999",
		"identity_card_blob_id": "blob_v1_12312",
		"avatar_blob_id":        "avatar_v1_445",
		"first_name":            "Hương",
		"last_name":             "Phạm Mai",
		"gender":                "Female",
		"date_of_birth":         "1999-01-01",
		"phone_number":          "+84934567890",
		"email":                 "huong.pham@example.com",
		"uploaded_at":           "1698398100000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft5 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[4],
		},
		"owner":                 sampleAdminIds[4],
		"identity_code":         "052090000555",
		"identity_card_blob_id": "blob_v1_88776",
		"avatar_blob_id":        "avatar_v1_556",
		"first_name":            "Tuấn",
		"last_name":             "Ngô Anh",
		"gender":                "Male",
		"date_of_birth":         "1990-03-30",
		"phone_number":          "+84945678901",
		"email":                 "tuan.ngo@example.com",
		"uploaded_at":           "1698398400000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft6 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[5],
		},
		"owner":                 sampleAdminIds[5],
		"identity_code":         "001085000111",
		"identity_card_blob_id": "blob_v1_66554",
		"avatar_blob_id":        "avatar_v1_667",
		"first_name":            "Thảo",
		"last_name":             "Vũ Thu",
		"gender":                "Female",
		"date_of_birth":         "1985-06-06",
		"phone_number":          "+84956789012",
		"email":                 "thao.vu@example.com",
		"uploaded_at":           "1698398700000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft7 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[6],
		},
		"owner":                 sampleAdminIds[6],
		"identity_code":         "038096000333",
		"identity_card_blob_id": "blob_v1_55443",
		"avatar_blob_id":        "avatar_v1_889",
		"first_name":            "Ngọc",
		"last_name":             "Bùi Bích",
		"gender":                "Female",
		"date_of_birth":         "1996-02-14",
		"phone_number":          "+84978901234",
		"email":                 "ngoc.bui@example.com",
		"uploaded_at":           "1698399300000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft8 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[7],
		},
		"owner":                 sampleAdminIds[7],
		"identity_code":         "001099000999",
		"identity_card_blob_id": "blob_v1_12312",
		"avatar_blob_id":        "avatar_v1_445",
		"first_name":            "Hương",
		"last_name":             "Phạm Mai",
		"gender":                "Female",
		"date_of_birth":         "1999-01-01",
		"phone_number":          "+84934567890",
		"email":                 "huong.pham@example.com",
		"uploaded_at":           "1698398100000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft9 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[8],
		},
		"owner":                 sampleAdminIds[8],
		"identity_code":         "060091000444",
		"identity_card_blob_id": "blob_v1_22110",
		"avatar_blob_id":        "avatar_v1_990",
		"first_name":            "Dũng",
		"last_name":             "Đỗ Tiến",
		"gender":                "Male",
		"date_of_birth":         "1991-09-09",
		"phone_number":          "+84989012345",
		"email":                 "dung.do@example.com",
		"uploaded_at":           "1698399600000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
	sampleJsonAdminNft10 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleAdminNfts[9],
		},
		"owner":                 sampleAdminIds[9],
		"identity_code":         "001094000555",
		"identity_card_blob_id": "blob_v1_00998",
		"avatar_blob_id":        "avatar_v1_001",
		"first_name":            "Trang",
		"last_name":             "Phan Huyền",
		"gender":                "Female",
		"date_of_birth":         "1994-04-04",
		"phone_number":          "+84990123456",
		"email":                 "trang.phan@example.com",
		"uploaded_at":           "1698399900000",
		"name":                  "RaiseChild Admin NFT",
		"url":                   "https://nft-storage.com",
	}
)

// Local Leader NFT
var (
	sampleJsonLeaderNft1 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[0],
		},
		"owner":                 sampleLocalLeaderIds[0],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[0],
		"identity_code":         "STF-001",
		"identity_card_blob_id": "bafybeic...key1",
		"avatar_blob_id":        "bafybeid...img1",
		"first_name":            "Minh",
		"last_name":             "Nguyen",
		"gender":                "Male",
		"date_of_birth":         "1992-04-12",
		"phone_number":          "+84912345678",
		"email":                 "minh.nguyen@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft2 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[1],
		},
		"owner":                 sampleLocalLeaderIds[1],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[1],
		"identity_code":         "STF-002",
		"identity_card_blob_id": "bafybeic...key2",
		"avatar_blob_id":        "bafybeid...img2",
		"first_name":            "Linh",
		"last_name":             "Tran",
		"gender":                "Female",
		"date_of_birth":         "1995-08-20",
		"phone_number":          "+84987654321",
		"email":                 "linh.tran@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft3 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[2],
		},
		"owner":                 sampleLocalLeaderIds[2],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[2],
		"identity_code":         "STF-003",
		"identity_card_blob_id": "bafybeic...key3",
		"avatar_blob_id":        "bafybeid...img3",
		"first_name":            "Hoang",
		"last_name":             "Pham",
		"gender":                "Male",
		"date_of_birth":         "1988-11-30",
		"phone_number":          "+84901122334",
		"email":                 "hoang.pham@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft4 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[3],
		},
		"owner":                 sampleLocalLeaderIds[3],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[3],
		"identity_code":         "STF-004",
		"identity_card_blob_id": "bafybeic...key4",
		"avatar_blob_id":        "bafybeid...img4",
		"first_name":            "An",
		"last_name":             "Le",
		"gender":                "Female",
		"date_of_birth":         "1994-01-15",
		"phone_number":          "+84933445566",
		"email":                 "an.le@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft5 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[4],
		},
		"owner":                 sampleLocalLeaderIds[4],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[4],
		"identity_code":         "STF-005",
		"identity_card_blob_id": "bafybeic...key5",
		"avatar_blob_id":        "bafybeid...img5",
		"first_name":            "Bach",
		"last_name":             "Vu",
		"gender":                "Male",
		"date_of_birth":         "1990-06-25",
		"phone_number":          "+84944556677",
		"email":                 "bach.vu@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft6 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[5],
		},
		"owner":                 sampleLocalLeaderIds[5],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[5],
		"identity_code":         "STF-006",
		"identity_card_blob_id": "bafybeic...key6",
		"avatar_blob_id":        "bafybeid...img6",
		"first_name":            "Thao",
		"last_name":             "Do",
		"gender":                "Female",
		"date_of_birth":         "1993-09-05",
		"phone_number":          "+84955667788",
		"email":                 "thao.do@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft7 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[6],
		},
		"owner":                 sampleLocalLeaderIds[6],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[6],
		"identity_code":         "STF-007",
		"identity_card_blob_id": "bafybeic...key7",
		"avatar_blob_id":        "bafybeid...img7",
		"first_name":            "Dung",
		"last_name":             "Bui",
		"gender":                "Male",
		"date_of_birth":         "1987-03-18",
		"phone_number":          "+84966778899",
		"email":                 "dung.bui@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft8 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[7],
		},
		"owner":                 sampleLocalLeaderIds[7],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[7],
		"identity_code":         "STF-008",
		"identity_card_blob_id": "bafybeic...key8",
		"avatar_blob_id":        "bafybeid...img8",
		"first_name":            "Nhi",
		"last_name":             "Phan",
		"gender":                "Female",
		"date_of_birth":         "1996-07-12",
		"phone_number":          "+84977889900",
		"email":                 "nhi.phan@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft9 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[8],
		},
		"owner":                 sampleLocalLeaderIds[8],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[8],
		"identity_code":         "STF-009",
		"identity_card_blob_id": "bafybeic...key9",
		"avatar_blob_id":        "bafybeid...img9",
		"first_name":            "Kien",
		"last_name":             "Dang",
		"gender":                "Male",
		"date_of_birth":         "1991-12-01",
		"phone_number":          "+84988990011",
		"email":                 "kien.dang@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonLeaderNft10 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleLocalLeaderNfts[9],
		},
		"owner":                 sampleLocalLeaderIds[9],
		"role":                  "Local Leader",
		"region":                sampleLocalRegions[9],
		"identity_code":         "STF-010",
		"identity_card_blob_id": "bafybeic...key10",
		"avatar_blob_id":        "bafybeid...img10",
		"first_name":            "Mai",
		"last_name":             "Hoang",
		"gender":                "Female",
		"date_of_birth":         "1994-05-22",
		"phone_number":          "+84999001122",
		"email":                 "mai.hoang@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Local Leader NFT",
		"url":                   "https://arweave.net",
	}
)

// Volunteer NFT
var (
	sampleJsonVolunteerNft1 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[0],
		},
		"owner":                 sampleVolunteerIds[0],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[0],
		"identity_code":         "STF-001",
		"identity_card_blob_id": "bafybeic...key1",
		"avatar_blob_id":        "bafybeid...img1",
		"first_name":            "Minh",
		"last_name":             "Nguyen",
		"gender":                "Male",
		"date_of_birth":         "1992-04-12",
		"phone_number":          "+84912345678",
		"email":                 "minh.nguyen@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft2 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[1],
		},
		"owner":                 sampleVolunteerIds[1],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[1],
		"identity_code":         "STF-002",
		"identity_card_blob_id": "bafybeic...key2",
		"avatar_blob_id":        "bafybeid...img2",
		"first_name":            "Linh",
		"last_name":             "Tran",
		"gender":                "Female",
		"date_of_birth":         "1995-08-20",
		"phone_number":          "+84987654321",
		"email":                 "linh.tran@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft3 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[2],
		},
		"owner":                 sampleVolunteerIds[2],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[2],
		"identity_code":         "STF-003",
		"identity_card_blob_id": "bafybeic...key3",
		"avatar_blob_id":        "bafybeid...img3",
		"first_name":            "Hoang",
		"last_name":             "Pham",
		"gender":                "Male",
		"date_of_birth":         "1988-11-30",
		"phone_number":          "+84901122334",
		"email":                 "hoang.pham@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft4 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[3],
		},
		"owner":                 sampleVolunteerIds[3],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[3],
		"identity_code":         "STF-004",
		"identity_card_blob_id": "bafybeic...key4",
		"avatar_blob_id":        "bafybeid...img4",
		"first_name":            "An",
		"last_name":             "Le",
		"gender":                "Female",
		"date_of_birth":         "1994-01-15",
		"phone_number":          "+84933445566",
		"email":                 "an.le@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft5 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[4],
		},
		"owner":                 sampleVolunteerIds[4],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[4],
		"identity_code":         "STF-005",
		"identity_card_blob_id": "bafybeic...key5",
		"avatar_blob_id":        "bafybeid...img5",
		"first_name":            "Bach",
		"last_name":             "Vu",
		"gender":                "Male",
		"date_of_birth":         "1990-06-25",
		"phone_number":          "+84944556677",
		"email":                 "bach.vu@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft6 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[5],
		},
		"owner":                 sampleVolunteerIds[5],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[5],
		"identity_code":         "STF-006",
		"identity_card_blob_id": "bafybeic...key6",
		"avatar_blob_id":        "bafybeid...img6",
		"first_name":            "Thao",
		"last_name":             "Do",
		"gender":                "Female",
		"date_of_birth":         "1993-09-05",
		"phone_number":          "+84955667788",
		"email":                 "thao.do@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft7 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[6],
		},
		"owner":                 sampleVolunteerIds[6],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[6],
		"identity_code":         "STF-007",
		"identity_card_blob_id": "bafybeic...key7",
		"avatar_blob_id":        "bafybeid...img7",
		"first_name":            "Dung",
		"last_name":             "Bui",
		"gender":                "Male",
		"date_of_birth":         "1987-03-18",
		"phone_number":          "+84966778899",
		"email":                 "dung.bui@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft8 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[7],
		},
		"owner":                 sampleVolunteerIds[7],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[7],
		"identity_code":         "STF-008",
		"identity_card_blob_id": "bafybeic...key8",
		"avatar_blob_id":        "bafybeid...img8",
		"first_name":            "Nhi",
		"last_name":             "Phan",
		"gender":                "Female",
		"date_of_birth":         "1996-07-12",
		"phone_number":          "+84977889900",
		"email":                 "nhi.phan@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft9 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[8],
		},
		"owner":                 sampleVolunteerIds[8],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[8],
		"identity_code":         "STF-009",
		"identity_card_blob_id": "bafybeic...key9",
		"avatar_blob_id":        "bafybeid...img9",
		"first_name":            "Kien",
		"last_name":             "Dang",
		"gender":                "Male",
		"date_of_birth":         "1991-12-01",
		"phone_number":          "+84988990011",
		"email":                 "kien.dang@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
	sampleJsonVolunteerNft10 = map[string]interface{}{
		"id": map[string]string{
			"id": sampleVolunteerNfts[9],
		},
		"owner":                 sampleVolunteerIds[9],
		"role":                  "Volunteer",
		"region":                sampleLocalRegions[9],
		"identity_code":         "STF-010",
		"identity_card_blob_id": "bafybeic...key10",
		"avatar_blob_id":        "bafybeid...img10",
		"first_name":            "Mai",
		"last_name":             "Hoang",
		"gender":                "Female",
		"date_of_birth":         "1994-05-22",
		"phone_number":          "+84999001122",
		"email":                 "mai.hoang@corp.it",
		"uploaded_at":           "1709283600000",
		"name":                  "RaiseChild Volunteer NFT",
		"url":                   "https://arweave.net",
	}
)

const (
	sampleManageObjId string = ""
)

var (
	sampleAdminIds              = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleAdminNfts             = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleChildIds              = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleVolunteerIds          = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleVolunteerNfts         = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleLocalLeaderIds        = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleLocalLeaderNfts       = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleLocalRegions          = []string{shared.HA_NOI_REGION, shared.HO_CHI_MINH_REGION, "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleChildrenCenters       = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleCenterConfirmStatuses = []bool{true, true, true, true, false, false, true, true, true, true}
	sampleCreatedCenters        = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleDonorIds              = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	sampleDonorNfts             = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
)

var (
	sampleUploadedBankProfile = entities.BankProfile{
		Owner:     "Owner",
		BankOrg:   "Bank Owner",
		OwnerName: "Owner Name",
	}
)

func getFullJsonAdminNfts() []map[string]interface{} {
	return []map[string]interface{}{
		sampleJsonAdminNft1,
		sampleJsonAdminNft2,
		sampleJsonAdminNft3,
		sampleJsonAdminNft4,
		sampleJsonAdminNft5,
		sampleJsonAdminNft6,
		sampleJsonAdminNft7,
		sampleJsonAdminNft8,
		sampleJsonAdminNft9,
		sampleJsonAdminNft10,
	}
}

func getFoundJsonAdminNftsWithKeyWord() ([]map[string]interface{}, string) {
	return []map[string]interface{}{
		sampleJsonAdminNft1,
		sampleJsonAdminNft3,
		sampleJsonAdminNft4,
		sampleJsonAdminNft5,
		sampleJsonAdminNft7,
		sampleJsonAdminNft8,
		sampleJsonAdminNft9,
		sampleJsonAdminNft10,
	}, "ng"
}
