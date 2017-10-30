package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.Argument;
import com.marshmallow.anwork.app.cli.Command;
import com.marshmallow.anwork.app.cli.Flag;
import com.marshmallow.anwork.app.cli.List;
import com.marshmallow.anwork.app.cli.Visitor;

import java.util.ArrayList;

/**
 * This is a {@link Visitor} that is used to test the CLI visitor functionality.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class TestCliVisitor implements Visitor {

  private java.util.List<String> visitedShortFlags = new ArrayList<String>();
  private java.util.List<String> visitedShortFlagsWithParameters = new ArrayList<String>();
  private java.util.List<String> visitedLongFlags = new ArrayList<String>();
  private java.util.List<String> visitedLongFlagsWithParameters = new ArrayList<String>();
  private java.util.List<String> visitedCommands = new ArrayList<String>();
  private java.util.List<String> visitedCommandArguments = new ArrayList<String>();
  private java.util.List<String> visitedLists = new ArrayList<String>();
  private java.util.List<String> leftLists = new ArrayList<String>();

  /**
   * Get the short flags that were visited by this {@link Visitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link Visitor}
   */
  public String[] getVisitedShortFlags() {
    return visitedShortFlags.toArray(new String[0]);
  }

  /**
   * Get the short flags with parameters that were visited by this {@link Visitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link Visitor}
   */
  public String[] getVisitedShortFlagsWithParameters() {
    return visitedShortFlagsWithParameters.toArray(new String[0]);
  }

  /**
   * Get the long flags that were visited by this {@link Visitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link Visitor}
   */
  public String[] getVisitedLongFlags() {
    return visitedLongFlags.toArray(new String[0]);
  }

  /**
   * Get the long flags with parameters that were visited by this {@link Visitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link Visitor}
   */
  public String[] getVisitedLongFlagsWithParameters() {
    return visitedLongFlagsWithParameters.toArray(new String[0]);
  }

  /**
   * Get the commands that were visited by this {@link Visitor}. Note that the order of the
   * commands may not match the order in which they were added to a list!
   *
   * @return The commands that were visited by this {@link Visitor}
   */
  public String[] getVisitedCommands() {
    return visitedCommands.toArray(new String[0]);
  }

  /**
   * Get the arguments that were visited by this {@link Visitor}. The arguments appear in the order
   * in which they were added to each command. The order of commands follows the same pattern as
   * {@link #getVisitedCommands()}.
   *
   * @return The arguments that were visited by this {@link Visitor}
   */
  public String[] getVisitedCommandArguments() {
    return visitedCommandArguments.toArray(new String[0]);
  }

  /**
   * Get the lists that were visited by this {@link Visitor}. Note that the order of the lists
   * may not match the order in which they were added to a list!
   *
   * @return The lists that were visited by this {@link Visitor}
   */
  public String[] getVisitedLists() {
    return visitedLists.toArray(new String[0]);
  }

  /**
   * Get the lists that were left by this {@link Visitor}. Note that the order of the lists
   * may not match the order in which they were added to a list!
   *
   * @return The lists that were left by this {@link Visitor}
   */
  public String[] getLeftLists() {
    return leftLists.toArray(new String[0]);
  }

  @Override
  public void visitFlag(Flag flag) {
    if (flag.hasLongFlag()) {
      if (flag.hasArgument()) {
        visitedLongFlagsWithParameters.add(flag.getLongFlag());
      } else {
        visitedLongFlags.add(flag.getLongFlag());
      }
    } else {
      if (flag.hasArgument()) {
        visitedShortFlagsWithParameters.add(flag.getShortFlag());
      } else {
        visitedShortFlags.add(flag.getShortFlag());
      }
    }
  }

  @Override
  public void visitList(List list) {
    visitedLists.add(list.getName());
  }

  @Override
  public void leaveList(List list) {
    leftLists.add(list.getName());
  }

  @Override
  public void visitCommand(Command command) {
    visitedCommands.add(command.getName());
    for (Argument argument : command.getArguments()) {
      visitedCommandArguments.add(argument.getName());
    }
  }
}
