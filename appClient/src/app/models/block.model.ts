export interface BlockTransaction {
    type: string
    txId: string
    validationCode: string
}

export interface BlockInfo {
    blockNumber: number
    channelId: string
    source: string
    noOfTransactions: string
    transactions: BlockTransaction[]
}