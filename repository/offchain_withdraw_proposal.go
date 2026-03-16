package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"
)

type offChainWithdrawProposalRepo struct {
	db        *sql.DB
	errLogger *log.Logger
}

const (
	offchain_withdraw_proposal_table string = "withdraw_proposals"
)

func InitializeOffChainWithdrawProposalRepository(db *sql.DB, errLogger *log.Logger) repository.IOffChainWithdrawProposalRepository {
	return &offChainWithdrawProposalRepo{
		db:        db,
		errLogger: errLogger,
	}
}

// CreateOffChainWithdrawProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalRepo) CreateOffChainWithdrawProposal(donation entities.OffChainWithdrawProposal, ctx context.Context) error {
	var query string = "INSERT INTO " + offchain_withdraw_proposal_table +
		" (id, purpose, proposal_id, target, local_pool_id, created_at) " +
		"values ($1, $2, $3, $4, $5)"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_WITHDRAW_PROPOSAL_REPOSITORY) + "CreateOffChainWithdrawProposal - "

	if _, err := o.db.ExecContext(ctx, query, donation.ID, donation.Purpose, donation.ProposalID, donation.Target, donation.LocalPoolID, donation.CreatedAt); err != nil {

		o.errLogger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetOffChainWithdrawProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalRepo) GetOffChainWithdrawProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error) {
	var query string = "SELECT * FROM " + offchain_withdraw_proposal_table + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_WITHDRAW_PROPOSAL_REPOSITORY) + "GetOffChainWithdrawProposal - "

	var res entities.OffChainWithdrawProposal
	if err := o.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Purpose, &res.ProposalID, &res.Target, &res.LocalPoolID, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		o.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetOffChainWithdrawProposalByProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalRepo) GetOffChainWithdrawProposalByProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error) {
	var query string = "SELECT * FROM " + offchain_withdraw_proposal_table + " WHERE proposal_id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_WITHDRAW_PROPOSAL_REPOSITORY) + "GetOffChainWithdrawProposalByProposal - "

	var res entities.OffChainWithdrawProposal
	if err := o.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Purpose, &res.ProposalID, &res.Target, &res.LocalPoolID, &res.CreatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		o.errLogger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// SetOnChainProposalIdAfterExecuteTx implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalRepo) SetOnChainProposalIdAfterExecuteTx(id string, proposalId string, ctx context.Context) error {
	var query string = "UPDATE " + offchain_withdraw_proposal_table + " SET proposal_id = $1 WHERE id = $2"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, shared.OFFCHAIN_WITHDRAW_PROPOSAL_REPOSITORY) + "SetOnChainProposalIdAfterExecuteTx - "
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := o.db.ExecContext(ctx, query, proposalId, id)
	if err != nil {
		o.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		o.errLogger.Println(errLogMsg + err.Error())
		return internalErr
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, offchain_withdraw_proposal_table))
	}

	return nil
}
