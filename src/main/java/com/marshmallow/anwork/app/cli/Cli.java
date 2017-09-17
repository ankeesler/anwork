package com.marshmallow.anwork.app.cli;

/**
 * This class represents the entry point for the CLI tree.
 *
 * <p>The class is meant to be used in the following way.
 * <pre>
 *   Cli cli = new Cli();
 *   CliNode root = cli.getRoot();
 *   root.addShortFlag("f",
 *                     "This is a description", new CliAction() { ... });
 *   root.addLongFlag("v",
 *                    "verbose",
 *                    "This is a description",
 *                    new CliAction() { ... });
 *   root.addLongFlagWithParameter("o",
 *                                 "output",
 *                                 "This is a description",
 *                                 "location",
 *                                 new CliAction() { ... });
 *
 *   CliNode tunaCommand = root.addCommand("tuna",
 *                                         "This is a tuna command",
 *                                         new CliAction() { ... });
 *   tunaCommand.addShortFlag("a", "This is a description", new CliAction() { ... });
 *
 *   CliNode fishList = root.addList("fish", "This is the fish command list");
 *   CliNode fishMarlinCommand = fishList.addCommand("marlin",
 *                                                   "This is the marlin command",
 *                                                   new CliAction() { ... });
 *   ...
 *   root.parse(args);
 * </pre>
 *
 * <p>The above would result in the following command line API.
 * <pre>
 *   root-command [-f] [-v|--verbose] [-o|--output location] tuna [-a]
 *   root-command [-f] [-v|--verbose] [-o|--output location] fish marlin
 * </pre>
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class Cli {

  private final CliNodeImpl root;

  /**
   * Create a new command line interface.
   *
   * @param name The name for the command line interface
   * @param description The description of the command line interface
   */
  public Cli(String name, String description) {
    root = CliNodeImpl.makeRoot(name, description);
  }

  /**
   * Get the root CLI node for this CLI.
   *
   * @return The root CLI node
   */
  public CliList getRoot() {
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
}
