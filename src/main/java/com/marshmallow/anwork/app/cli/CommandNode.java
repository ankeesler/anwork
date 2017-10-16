package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link Node} that represents a {@link Command}.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
class CommandNode extends Node implements MutableCommand {

  CommandNode(String name, Action action) {
    super(name, action);
  }

  @Override
  protected boolean isList() {
    return false;
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
    super.setAction(action);
    return this;
  }

  @Override
  public Action getAction() {
    return super.getAction();
  }

  @Override
  protected void startVisit(Visitor visitor) {
    visitor.visitCommand(this);
  }

  @Override
  protected void endVisit(Visitor visitor) {
    // no-op
  }
}