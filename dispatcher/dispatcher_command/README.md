#### Usage


```go
//@Author: Nghiant5
// This funcion is a sample to register a handler for commandType
// This function you can init it in folder configuration of your service
// And you Dependency Inject your dispatcher and use it to register
// And Interface Handler used for handler -> and the handler sample is below this function
func InitCommandHandlerRegister(dispatcher dp.DispatcherHandlerSaga, executeCommand domain.ExecuteNapasPaymentCommand, emitCommand domain.ExecuteNapasPaymentEmitHandler) {

	dispatcher.RegisterHandler(string(command.TransactionReceivedCommandType), emitCommand.EmitSagaCreate)
	dispatcher.RegisterHandler(string(command.TransactionCreatedCommandType), executeCommand.ValidateAccount)
	dispatcher.RegisterHandler(string(command.CurrentAccountValidatedCommandType), executeCommand.FundTransferExecute)
	dispatcher.RegisterHandler(string(command.FundTransferResponseCommandType), emitCommand.EmitFundTransferResult)
	dispatcher.RegisterHandler(string(command.FundTransferSuccessCommandType), executeCommand.SendingNapas)
	dispatcher.RegisterHandler(string(command.NapasAdapterSuccessCommandType), emitCommand.EmitNapasResult)
	dispatcher.RegisterHandler(string(command.NapasExecutedCommandType), executeCommand.PublishResult)
}

type ExecuteNapasPaymentCommand interface {
	ValidateAccount(ctx context.Context, esEvent events.Event) error
	FundTransferExecute(ctx context.Context, esEvent events.Event) error
	SendingNapas(ctx context.Context, esEvent events.Event) error
	PublishResult(ctx context.Context, esEvent events.Event) error
}

type ExecuteNapasPaymentEmitHandler interface {
	EmitSagaCreated(ctx context.Context, esEvent events.Event) error
	EmitFundTransferResult(ctx context.Context, esEvent events.Event) error
	EmitNapasResult(ctx context.Context, esEvent events.Event) error

}