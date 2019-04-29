package main

type Permission string

const PermissionAll Permission = "PermissionAll"
const PermissionRead Permission = "PermissionRead"
const PermissionCreateWallet Permission = "PermissionCreateWallet"
const PermissionCreateOrder Permission = "PermissionCreateOrder"
const PermissionCancelOrder Permission = "PermissionCancelOrder"
const PermissionTokenBurn Permission = "PermissionTokenBurn"
const PermissionDeposit Permission = "PermissionDeposit"
const PermissionFreezeToken Permission = "PermissionFreezeToken"
const PermissionIssueToken Permission = "PermissionIssueToken"
const PermissionListPair Permission = "PermissionListPair"
const PermissionMintToken Permission = "PermissionMintToken"
const PermissionSendToken Permission = "PermissionSendToken"
const PermissionSubmitProposal Permission = "PermissionSubmitProposal"
const PermissionUnfreezeToken Permission = "PermissionUnfreezeToken"
const PermissionVoteProposal Permission = "PermissionVoteProposal"
