package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link ListOrCommandImpl} that represents a {@link Command}.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
class CommandImpl extends ListOrCommandImpl implements MutableCommand {

  private Action action;

  CommandImpl(ListImpl parent, String name, Action action) {
    super(parent, name);
    this.action = action;
  }

  @Override
  public MutableCommand setName(String name) {
    super.setName(name);
    return this;
  }

  @Override
  public MutableCommand setDescription(String description) {
    super.setDescription(description);
    return this;
  }

  @Override
  public MutableCommand setAction(Action action) {
    this.action = action;
    return this;
  }

  @Override
  public Action getAction() {
    return action;
  }

  // This is package-private so that it can be called from List parsing.
  void parse(String[] args, int index, ParseContext context) {
    context.setActiveNode(this);
    while (index < args.length) {
      if (isFlag(args[index])) {
        index = parseFlag(args, index, context);
      } else {
        context.addParameter(args[index]);
        index += 1;
      }
    }

    validateContext(context);
    runActionFromContext(context);
  }

  private void validateContext(ParseContext context) {

  }

  private void runActionFromContext(ParseContext context) {
    action.run(context.getFlagValues(), context.getParameters());
  }

  // This is package-private so that it can be called from List visitation.
  void visit(Visitor visitor) {
    visitFlags(visitor);
    visitor.visitCommand(this);
  }
}