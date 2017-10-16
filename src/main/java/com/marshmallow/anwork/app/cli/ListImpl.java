package com.marshmallow.anwork.app.cli;

import java.util.stream.Stream;

/**
 * This is a {@link ListOrCommandImpl} that represents a {@link List}.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
class ListImpl extends ListOrCommandImpl implements MutableList {

  /**
   * This is a simple {@link Action} that prints some usage information on a {@link List}.
   *
   * <p>
   * Created Oct 16, 2017
   * </p>
   *
   * @author Andrew
   */
  private static class ListAction implements Action {

    private final List list;

    public ListAction(List list) {
      this.list = list;
    }

    @Override
    public void run(ArgumentValues flags, String[] parameters) {
      StringBuilder builder = new StringBuilder();
      generateUsage(list, builder);
      System.out.println(builder.toString());
    }

    private static void generateUsage(List list, StringBuilder builder) {
      TextDocumentationGenerator.writeList(new List[] { list }, builder);
      if (list.getFlags().length > 0) {
        TextDocumentationGenerator.writeFlags(list.getFlags(), builder);
      }
      if (list.getCommands().length > 0) {
        TextDocumentationGenerator.writeCommands(list.getCommands(), builder);
      }
      Stream.of(list.getLists()).forEach((l) -> generateUsage(l, builder));
    }
  }

  ListImpl(String name) {
    super(name, null);
    setAction(new ListAction(this));
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
    ListImpl list = new ListImpl(name);
    addChild(list);
    return list;
  }

  @Override
  public List[] getLists() {
    return super.getLists();
  }

  @Override
  public MutableCommand addCommand(String name, Action action) {
    CommandImpl command = new CommandImpl(name, action);
    addChild(command);
    return command;
  }

  @Override
  public Command[] getCommands() {
    return super.getCommands();
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