package com.marshmallow.anwork.app.cli;

import java.util.ArrayList;
import java.util.List;

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
  private final List<Argument> arguments = new ArrayList<Argument>();

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

  @Override
  public <T> MutableArgument addArgument(String name, ArgumentType<T> type) {
    ArgumentImpl argument = new ArgumentImpl(name, type);
    arguments.add(argument);
    return argument;
  }

  @Override
  public Argument[] getArguments() {
    return arguments.toArray(new Argument[0]);
  }

  // This is package-private so that it can be called from List parsing.
  void parse(String[] args, int index, ArgumentValues flagValues) {
    List<String> arguments = new ArrayList<String>();
    while (index < args.length) {
      if (isFlag(args[index])) {
        index = parseFlag(args, index, flagValues);
      } else {
        arguments.add(args[index]);
        index += 1;
      }
    }

    ArgumentValues argumentValues = makeArgumentValues(arguments.toArray(new String[0]));
    action.run(flagValues, argumentValues);
  }

  private ArgumentValues makeArgumentValues(String[] actualArguments) {
    Argument[] expectedArguments = getArguments();
    validateArgumentCount(expectedArguments, actualArguments);

    ArgumentValues argumentValues = new ArgumentValues();
    for (int i = 0; i < expectedArguments.length; i++) {
      Argument expectedArgument = expectedArguments[i];
      String actualArgument = actualArguments[i];
      Object argumentValue = validateArgumentType(expectedArgument, actualArgument);
      argumentValues.addValue(expectedArgument.getName(), argumentValue);
    }
    return argumentValues;
  }

  private void validateArgumentCount(Argument[] expectedArguments, String[] actualArguments) {
    if (actualArguments.length != expectedArguments.length) {
      StringBuilder builder = new StringBuilder("Incorrect number of arguments passed to "
                                                + "'" + getName() + "' command. Expected "
                                                + expectedArguments.length);
      for (Argument expectedArgument : expectedArguments) {
        builder.append(' ');
        builder.append('<');
        builder.append(expectedArgument.getType().getName());
        builder.append(' ');
        builder.append(expectedArgument.getName());
        builder.append('>');
      }
      throw new IllegalArgumentException(builder.toString());
    }
  }

  private Object validateArgumentType(Argument expectedArgument, String actualArgument) {
    try {
      return expectedArgument.getType().convert(actualArgument);
    } catch (IllegalArgumentException iae) {
      throw new IllegalArgumentException(("Cannot convert argument '" + actualArgument + "' "
                                          + " to a " + expectedArgument.getType().getName()),
                                         iae);
    }
  }

  // This is package-private so that it can be called from List visitation.
  void visit(Visitor visitor) {
    visitFlags(visitor);
    visitor.visitCommand(this);
  }
}