package com.marshmallow.anwork.app.cli;

/**
 * This class represents the entry point for the CLI API.
 *
 * <p>The class is meant to be used in the following way.
 * <pre>
 *   Cli cli = new Cli("cli-app");
 *   MutableList rootList = cli.getRoot();
 *   rootList.addFlag("f")
 *           .setDescription("This is a description");
 *   rootList.addFlag("o")
 *           .setLongFlag("output")
 *           .setDescription("This is a description")
 *           .setArgument("location").setType(ArgumentType.STRING);
 *
 *   MutableCommand tunaCommand = root.addCommand("tuna", new CliAction() { ... });
 *   tunaCommand.setDescription("This is the tuna command");
 *   tunaCommand.addFlag("a").setDescription("This is a description");
 *
 *   MutableList fishList = root.addList("fish").setDescription(""This is the fish command list");
 *   fishList.addCommand("marlin", new CliAction() { ... }))
 *           .setDescription("This is the marlin command");
 *   ...
 *   root.parse(args);
 * </pre>
 *
 * <p>The above would result in the following command line API.
 * <pre>
 *   cli-app [-f] [-o|--output location] tuna [-a]
 *   cli-app [-f] [-o|--output location] fish marlin
 * </pre>
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class Cli {

  private final ListImpl root;

  /**
   * Create a new command line interface (i.e., a {@link Cli}).
   *
   * @param name The name for the command line interface
   */
  public Cli(String name) {
    root = new ListImpl(name);
  }

  /**
   * Get the root {@link List} for this CLI, in editable form.
   *
   * @return The root CLI {@link List} that can be edited, i.e., a {@link MutableList}
   */
  public MutableList getRoot() {
    return root;
  }

  /**
   * Parse the provided command line arguments and run the necessary actions.
   *
   * @param args The command line arguments
   * @throws IllegalArgumentException If the command line arguments are bad
   */
  public void parse(String[] args) throws IllegalArgumentException {
    root.parse(args);
  }

  /**
   * Get the usage information for this CLI node.
   *
   * @return The usage information for this CLI node
   */
  public String getUsage() {
    return root.getUsage();
  }

  /**
   * Visit the CLI tree with a {@link Visitor}.
   *
   * <p>Per {@link List} in the {@link Cli}, the visitation will happen in this order.
   *   <ol>
   *     <li>Flags    (see {@link Flag})</li>
   *     <li>Commands (see {@link Command})</li>
   *     <li>Lists    (see {@link List})</li>
   *   </ol>
   * Each of the groups above will be sorted so that each visitation sequence is deterministic.
   * When lists are visited, they are visited in a depth-first manner.
   *
   * @param visitor The {@link Visitor} with which to visit the CLI tree nodes
   */
  public void visit(Visitor visitor) {
    root.visit(visitor);
  }
}
