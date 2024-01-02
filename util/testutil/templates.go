package testutil

const (
	AuthzGrantTxTemplates                   = `{"body":{"messages":[{"@type":"/cosmos.authz.v1beta1.MsgGrant","granter":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","grantee":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","grant":{"authorization":{"@type":"/cosmos.bank.v1beta1.SendAuthorization","spend_limit":[{"denom":"axpla","amount":"1000"}]},"expiration":"1970-01-01T00:00:00Z"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["CELJ72C7P1PYgUN3RRyy7D7MfG/vBQmwGhN8hjoPWng7GZlb3/wI4sW+DtBufvHN+Z8sIL2j9bOqujmkz9fijQA="]}`
	AuthzRevokeTxTemplates                  = `{"body":{"messages":[{"@type":"/cosmos.authz.v1beta1.MsgRevoke","granter":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","grantee":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","msg_type_url":"/cosmos.bank.v1beta1.MsgSend"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["BNn4HUUOczlDCgemxs7qka2dMDB1myJeX8mapBTgtR1439pxSEM0vwhflBtK6QhegRfd20+e4nM4p9zgeycA2gE="]}`
	AuthzExecTxTemplates                    = `{"body":{"messages":[{"@type":"/cosmos.authz.v1beta1.MsgExec","grantee":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","msgs":[{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","to_address":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","amount":[{"denom":"axpla","amount":"1000"}]}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["qI4LP4tlzYpzI5qvjfYgnGdcZe4ZjZvkjxIBjj9QVxhASTubQN14r01qC4IT+0TJVkJOOvVBdw1ObCzGBB5iKAA="]}`
	BankSendTxTemplates                     = `{"body":{"messages":[{"@type":"/cosmos.bank.v1beta1.MsgSend","from_address":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","to_address":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","amount":[{"denom":"axpla","amount":"1000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["8zHJ92yptbczT3dSXKgfD3DlZMCzPki8pBQXUKqcwrwKzLCedIbWfEnoHL2vMLIAB+PYbqhr+2fUZ5KlX4pKSgE="]}`
	CrisisInvariantBrokenTxTemplates        = `{"body":{"messages":[{"@type":"/cosmos.crisis.v1beta1.MsgVerifyInvariant","sender":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","invariant_module_name":"bank","invariant_route":"total-supply"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["/kTkdqHD6KlkdEgFurSLoLNZZaMW9aomo7RMplu5ooUcnpLcIYyLb80c2ocI31wdz9CU4H47oEciQxUUwhOL4AE="]}`
	DistFundCommunityPoolTxTemplates        = `{"body":{"messages":[{"@type":"/cosmos.distribution.v1beta1.MsgFundCommunityPool","amount":[{"denom":"axpla","amount":"1000"}],"depositor":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["lTw2DId49D8UUn0kKzO9TBYFLZNdEAPBVhELysJT+8JnVQloCF5g1mpTbfnwqlxXob+iIUvXl1bpMLyKPkV6QwE="]}`
	DistCommunityPoolSpendTxTemplates       = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.distribution.v1beta1.CommunityPoolSpendProposal","title":"community pool spend","description":"pay me","recipient":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","amount":[{"denom":"axpla","amount":"1000"}]},"initial_deposit":[{"denom":"axpla","amount":"1000"}],"proposer":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["hLlJhBW9kuGM+OnwqkwY4x4ixlv63iV2PpM6tzTziLokp2u3sVrMxnquPJQ1xTQG1QtpOIOp7mRFrGVQJByBMgA="]}`
	DistWithdrawRewardsTxTemplates          = `{"body":{"messages":[{"@type":"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward","delegator_address":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","validator_address":"xplavaloper1l03kma4vv9qcvhgcxf2ga0rnv7dqcumahy92l2"},{"@type":"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission","validator_address":"xplavaloper1l03kma4vv9qcvhgcxf2ga0rnv7dqcumahy92l2"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["sKk7YPLjuH8AusYsM0Cwb4rrnVtRDpr/K40yS2pANnwMEb7HsZ3pdyEiBaFHfLUnJX19Wd1ojIrobHsE9zOxzgE="]}`
	DistSetWithdrawAddrTxTemplates          = `{"body":{"messages":[{"@type":"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress","delegator_address":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","withdraw_address":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["qQi+d0w0pAlPYnkfmRvJ2uPm2JMnD4zXEQC7Gegy/yJhm1GANtscyqRZMkBGCFeEdfkfikFbdx+lsv3fcVFyoAA="]}`
	FeegrantFeegrantTxTemplates             = `{"body":{"messages":[{"@type":"/cosmos.feegrant.v1beta1.MsgGrantAllowance","granter":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","grantee":"xpla1h4x2jlnqkzq2k8wrfzvttl9p5gcffz4xe5cj2c","allowance":{"@type":"/cosmos.feegrant.v1beta1.BasicAllowance","spend_limit":[{"denom":"axpla","amount":"1000"}],"expiration":"2100-01-01T23:59:59Z"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["5+SEWutMwqw8mSfUhS5B0TwcaoRkwrEtOPC5g5Y3NAw4STm+Qy/XmAnsjHdZltkYt6hkK1AgAT1vWCr8lI8DbAA="]}`
	FeegrantRevokeFeegrantTxTemplates       = `{"body":{"messages":[{"@type":"/cosmos.feegrant.v1beta1.MsgRevokeAllowance","granter":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","grantee":"xpla1h4x2jlnqkzq2k8wrfzvttl9p5gcffz4xe5cj2c"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["kd9oYbFdlK3y6U/iTafh6xCR10v+R8h8Od95qTv6JhEIMcsi+XbHNYOAzEjGODSR95oJicUSHnuxAK45yT/MWgE="]}`
	GovSubmitProposalTxTemplates            = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.gov.v1beta1.TextProposal","title":"Test proposal","description":"Proposal description"},"initial_deposit":[{"denom":"axpla","amount":"1000"}],"proposer":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["5tnRLtVVb7EfdhLf2mOoz/FdiJa0lWipa7233pY4K0QunSgkAISO2KcjH/8zEVJ5v5VmjrlivqIYAQ870BZ4tgE="]}`
	GovDepositTxTemplates                   = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgDeposit","proposal_id":"1","depositor":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","amount":[{"denom":"axpla","amount":"1000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["qE9Ei7exI5ARVF7NT2wGYOLTjHgBQfC4m2+rxG1CzbMTimjAKprX14HhnqjwnP6vmtSOC5OhWyYh61pzfyOBEAE="]}`
	GovVoteTxTemplates                      = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgVote","proposal_id":"1","voter":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","option":"VOTE_OPTION_YES"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["+bNm10eijhVD6RpStn3RopnLj0sraFvbIF3VNk9lGRcU4I+nRDR3RdcuLz5LT5x8mOQ66BtoEG5SY9TgawGeCQA="]}`
	GovWeightedVoteTxTemplates              = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgVoteWeighted","proposal_id":"1","voter":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh","options":[{"option":"VOTE_OPTION_YES","weight":"0.600000000000000000"},{"option":"VOTE_OPTION_ABSTAIN","weight":"0.050000000000000000"},{"option":"VOTE_OPTION_NO","weight":"0.300000000000000000"},{"option":"VOTE_OPTION_NO_WITH_VETO","weight":"0.050000000000000000"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["2aWDylQmslN82wi1+9j5+SpM8BbFS5S7QB3juZbe+YttR1t/QaKft/9j4NDHhzXq/gAURLJArChZRYAIHggnXAA="]}`
	ParamsParamChangeTxTemplates            = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.params.v1beta1.ParameterChangeProposal","title":"Staking param change","description":"update max validators","changes":[{"subspace":"staking","key":"MaxValidators","value":"105"}]},"initial_deposit":[{"denom":"axpla","amount":"1000"}],"proposer":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["oM7bvJd89LQ4fFXy+wEIl+ATZBimv4IRDMLfRB8+5TxsS9qMeuSM1Ah5mDL8ZHqfWSbQ9VSlIkdrM81uRT/BkAA="]}`
	RewardFundFeeCollectorTxTemplates       = `{"body":{"messages":[{"@type":"/xpla.reward.v1beta1.MsgFundFeeCollector","amount":[{"denom":"axpla","amount":"1000"}],"depositor":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["yn0pvmUkTHOFa7vQniwnTGFSV3omwgEcsMJrgYEN5i8B743XGayuSSA0ls+5AwsMgtEKKlZgC6A+h3xDozlm/AA="]}`
	SlashingUnjailTxTemplates               = `{"body":{"messages":[{"@type":"/cosmos.slashing.v1beta1.MsgUnjail","validator_addr":"xplavaloper1l03kma4vv9qcvhgcxf2ga0rnv7dqcumahy92l2"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["hPMdSG0GJUGDRjQPHF6k+wxQtJUtH52jcXfcUMaxAcxV9Ar+6oCadkMckM1dqerJGM8Vhcw3wIylytPZpAbxKQE="]}`
	StakingEditValidatorTxTemplates         = `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgEditValidator","description":{"moniker":"moniker","identity":"identity","website":"website","security_contact":"securityContact","details":"details"},"validator_address":"xplavaloper1l8l7uju593qtu08uprtrly223dnpxlrvwmmmmg","commission_rate":null,"min_self_delegation":null}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["GKG3tx0/ggvp+s7GlzmAVm8X3uETLpvKufPWHeNWT78tZI8Je6a/3yTXR4UbypfTsw8OH5H+P9jkSjgMCo2p5AA="]}`
	StakingDelegateTxTemplates              = `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgDelegate","delegator_address":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","validator_address":"xplavaloper1l8l7uju593qtu08uprtrly223dnpxlrvwmmmmg","amount":{"denom":"axpla","amount":"1000"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["7Azy7e4tkF8t+SsIfnn/UtGR2E6mStR0ygTrELlj4RoZxmmJbrcqbIqoI/JYZb6lFVaiNrq3+/kCzaRuiMUCqQA="]}`
	StakingUnbondTxTemplates                = `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgUndelegate","delegator_address":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","validator_address":"xplavaloper1l8l7uju593qtu08uprtrly223dnpxlrvwmmmmg","amount":{"denom":"axpla","amount":"1000"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["cMdHHr0Y40q39hlWJF9H0Tbw10h9nIJuI6thpHRhDX5TFcjUoZDul00RS+6P76iSxLa+wPdOatZ1EG0PaHv3/AE="]}`
	StakingRedelegateTxTemplates            = `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgBeginRedelegate","delegator_address":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","validator_src_address":"xplavaloper1l8l7uju593qtu08uprtrly223dnpxlrvwmmmmg","validator_dst_address":"xplavaloper1l03kma4vv9qcvhgcxf2ga0rnv7dqcumahy92l2","amount":{"denom":"axpla","amount":"1000"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["Tj80GLrhSqsylMTIARlncL2GXj/W/RRg4NsGAOpxGuZBjqajjqrZvheWa44l2DBBPVm7iM1AOhFq1aSIJVYnrQA="]}`
	UpgradeSoftwareUpgradeTxTemplates       = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal","title":"Upgrade Title","description":"Upgrade Description","plan":{"name":"Upgrade Name","time":"0001-01-01T00:00:00Z","height":"6000","info":"{\"upgrade_info\":\"INFO\"}","upgraded_client_state":null}},"initial_deposit":[{"denom":"axpla","amount":"1000"}],"proposer":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["cobmhnGFWyVrcOLF+sVtZ+4ZZ6PrXcpCVdtLAFZ6vNse3d2szfgArPjATMvXKk7A15XBbuUn7TCizV6QMxvRkwE="]}`
	UpgradeCancelSoftwareUpgradeTxTemplates = `{"body":{"messages":[{"@type":"/cosmos.gov.v1beta1.MsgSubmitProposal","content":{"@type":"/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal","title":"1000","description":"Cancel software upgrade description"},"initial_deposit":[{"denom":"axpla","amount":"1000"}],"proposer":"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A5+rNG/0BpZEQGGZKq29JH4nvDnyHYmm1D+b5NzNC7bC"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["Wk26NO9tYdxFLpVp37/bgvF3MbCuaohq/9P7NAtUdZwQdSi9weQzsVh/4kLIpjxQDOUGMIL00akEc6KO5JpzogE="]}`
	WasmInstantiateContractTxTemplates      = `{"body":{"messages":[{"@type":"/cosmwasm.wasm.v1.MsgInstantiateContract","sender":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","admin":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","code_id":"1","label":"Contract instant","msg":{"owner":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54"},"funds":[{"denom":"axpla","amount":"10"}]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["0XFKN9YOzlJUmNM2fixrbuh9f02hT4oZL7Vcb9WN+rB0YgoQkDKbhrCU2LBu96iD+ckS/dxZA2/sUOBsiNUGrwE="]}`
	WasmExecuteContractTxTemplates          = `{"body":{"messages":[{"@type":"/cosmwasm.wasm.v1.MsgExecuteContract","sender":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","contract":"xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h","msg":{"execute_method":{"execute_key":"execute_test","execute_value":"execute_val"}},"funds":[]}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["I+BvuHu6Op74A99WZKAWJKSAUd04TaBClGOvqh5jpG0dOVYHy79hUJqor6bBRy9VvZptV5YhWQRj0LMcz0IaNQE="]}`
	WasmClearContractAdminTxTemplates       = `{"body":{"messages":[{"@type":"/cosmwasm.wasm.v1.MsgClearAdmin","sender":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","contract":"xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["oGpLH08YxSzAvqjezDosPE/0eiTYtNCwgJ4ihq+L1kA1Axu8DG68hwdQ2tnFWeOQce1qomVjhYy1MjhiiyUkUAA="]}`
	WasmSetContractAdminTxTemplates         = `{"body":{"messages":[{"@type":"/cosmwasm.wasm.v1.MsgUpdateAdmin","sender":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","new_admin":"","contract":"xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["eghZQ99M5Qe+vcVPF4NZ7mnHv76h2m03Gm37Wolr8NxpbsSUz6+663tOviuz5/OTiuiN2HzGdAxb6L4rRcS5OgE="]}`
	WasmMigrateTxTemplates                  = `{"body":{"messages":[{"@type":"/cosmwasm.wasm.v1.MsgMigrateContract","sender":"xpla1l8l7uju593qtu08uprtrly223dnpxlrvlxcp54","contract":"xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h","code_id":"2","msg":{}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9B1KwYOQUjFakc7Hhgbf1K/TldjpMWUvD5vWIMnHbF4"},"mode_info":{"single":{"mode":"SIGN_MODE_DIRECT"}},"sequence":"0"}],"fee":{"amount":[{"denom":"axpla","amount":"1275000000000000000"}],"gas_limit":"250000","payer":"","granter":""}},"signatures":["+8pSXcEzEOqmPbP3OxpoPITbo5AV8SjQeV7mCgRvqqEpIK1gsOOYJjNUDl3FNgfJMuasj9E3v4AAGloG+qwj0wA="]}`
)
