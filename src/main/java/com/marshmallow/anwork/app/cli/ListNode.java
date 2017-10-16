package com.marshmallow.anwork.app.cli;

/**
 * This is a {@link Node} that represents a {@link List}.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
class ListNode extends Node implements MutableList {

  ListNode(String name) {
    super(name, null);
    setAction(new Action() {
      @Override
      public void run(ArgumentValues flags, String[] arguments) {
        System.out.println(getUsage());
      }
    });
  }

  @Override
  protected boolean isList() {
    return true;
  }

  @Override
  public MutableList setName(String name) {
    super.setName(name);
    return this;
  }

  @Override
  public MutableList setDescription(String description) {
    super.setDescription(description);
    return this;
  }

  @Override
  public MutableList addList(String name) {
    ListNode list = new ListNode(name);
    addChild(list);
    return list;
  }

  @Override
  public List[] getLists() {
    return (List[])getChildren(true);
  }

  @Override
  public MutableCommand addCommand(String name, Action action) {
    CommandNode command = new CommandNode(name, action);
    addChild(command);
    return command;
  }

  @Override
  public Command[] getCommands() {
    return (Command[])getChildren(false);
  }

  @Override
  protected void startVisit(Visitor visitor) {
    visitor.visitList(this);
  }

  @Override
  protected void endVisit(Visitor visitor) {
    visitor.leaveList(this);
  }
}