package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.Command;
import com.marshmallow.anwork.app.cli.Flag;
import com.marshmallow.anwork.app.cli.List;
import com.marshmallow.anwork.app.cli.Visitor;

import java.util.ArrayList;

/**
 * This is a {@link Visitor} that records whether or not assorted optional data has been set on
 * objects in the CLI data structure.
 *
 * <p>
 * Created Oct 15, 2017
 * </p>
 *
 * @author Andrew
 */
public class OptionalDataCliVisitor implements Visitor {

  private final java.util.List<String> flagsWithDescriptions = new ArrayList<String>();
  private final java.util.List<String> commandsWithDescriptions = new ArrayList<String>();
  private final java.util.List<String> listsWithDescriptions = new ArrayList<String>();

  /**
   * Get an ordered list of the {@link Flag}'s (denoted by their short flags) that do have a
   * description, i.e., the {@link Flag}'s in the CLI where {@link Flag#hasDescription()} returns
   * <code>true</code>.
   *
   * @return The {@link Flag}'s in the CLI where {@link Flag#hasDescription()} returns
   *     <code>true</code>
   */
  public String[] getFlagsWithDescriptions() {
    return flagsWithDescriptions.toArray(new String[0]);
  }

  /**
   * Get an ordered list of the {@link Command}'s (denoted by their names) that do have a
   * description, i.e., the {@link Commands}'s in the CLI where {@link Command#hasDescription()}
   * returns <code>true</code>.
   *
   * @return The {@link Command}'s in the CLI where {@link Command#hasDescription()} returns
   *     <code>true</code>
   */
  public String[] getCommandsWithDescriptions() {
    return commandsWithDescriptions.toArray(new String[0]);
  }

  /**
   * Get an ordered list of the {@link List}'s (denoted by their names) that do have a
   * description, i.e., the {@link List}'s in the CLI where {@link List#hasDescription()}
   * returns <code>true</code>.
   *
   * @return The {@link List}'s in the CLI where {@link List#hasDescription()} returns
   *     <code>true</code>
   */
  public String[] getListsWithDescriptions() {
    return listsWithDescriptions.toArray(new String[0]);
  }

  @Override
  public void visitFlag(Flag flag) {
    if (flag.hasDescription()) {
      flagsWithDescriptions.add(flag.getShortFlag());
    }
  }

  @Override
  public void visitList(List list) {
    if (list.hasDescription()) {
      listsWithDescriptions.add(list.getName());
    }
  }

  @Override
  public void leaveList(List list) {
    // no-op
  }

  @Override
  public void visitCommand(Command command) {
    if (command.hasDescription()) {
      commandsWithDescriptions.add(command.getName());
    }
  }
}
