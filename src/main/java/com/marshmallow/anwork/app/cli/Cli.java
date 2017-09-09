package com.marshmallow.anwork.app.cli;

/**
 * This class represents the entry point for the CLI tree.
 *
 * The class is meant to be used in the following way.
 * <pre>
 *   Cli cli = new Cli();
 *   cli.addFlag(CliFlag.makeShortFlag("f", "This is a description", new CliAction() { ... });
 *   cli.addFlag(CliFlag.makeLongFlag("v", "verbose", "This is a description", new CliAction() { ... });
 *   cli.addFlag(CliFlag.makeLongFlagWithParameter("o", "output", "This is a description", "location", new CliAction() { ... });
 *
 *   CliNode tunaCommand = cli.addCommand("tuna", "This is a tuna command", new CliAction() { ... });
 *   cli.addFlag(CliFlag.makeShortFlag("a", "This is a description", new CliAction() { ... });
 *
 *   CliNode fishCommand = cli.addCommand("fish", "This is the fish command", new CliAction() { ... });
 *   CliNode fishMarlinCommand = cli.addCommand("marlin", "This is the marlin command", new CliAction() { ... });
 *   ...
 *   cli.parse(args);
 * </pre>
 *
 * The above would result in the following command line API.
 * <pre>
 *   *root* [-f] [-v|--verbose] [-o|--output location] tuna [-a]
 *   *root* [-f] [-v|--verbose] [-o|--output location] fish
 *   *root* [-f] [-v|--verbose] [-o|--output location] fish marlin 
 * </pre>
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class Cli implements CliNode {

  private static final String ROOT_NAME = "root";
  private static final String ROOT_DESCRIPTION = "root CLI node";

  private CliNode root
    = new CliNodeImpl(ROOT_NAME,
                      ROOT_DESCRIPTION,
                      (a) -> System.out.println(Cli.this.getUsage()));

  @Override
  public void addFlag(CliFlag flag) {
    root.addFlag(flag);
  }

  @Override
  public CliNode addCommand(String name, String description, CliAction action) {
    return root.addCommand(name, description, action);
  }

  @Override
  public void parse(String[] args) {
    root.parse(args);
  }

  @Override
  public String getUsage() {
    return root.getUsage();
  }

  @Override
  public String getName() {
    return root.getName();
  }

  @Override
  public String getDescription() {
    return root.getDescription();
  }
}
