package com.marshmallow.anwork.app.cli;

/**
 * This class represents the entry point for the CLI tree.
 *
 * The class is meant to be used in the following way.
 * <pre>
 *   Cli cli = new Cli("root-command");
 *   cli.addFlag(CliFlag.makeShortFlag("f", "This is a description", new CliAction() { ... });
 *   cli.addFlag(CliFlag.makeLongFlag("v", "verbose", "This is a description", new CliAction() { ... });
 *   cli.addFlag(CliFlag.makeLongFlagWithParameter("o", "output", "This is a description", "location", new CliAction() { ... });
 *
 *   CliNode tunaCommand = CliNode.makeCommand("tuna", "This is a tuna command", new CliAction() { ... });
 *   tunaCommand.addFlag(CliFlag.makeShortFlag("a", "This is a description", new CliAction() { ... });
 *   cli.addNode(tunaCommand);
 *
 *   CliNode fishList = CliNode.makeList("fish", "This is the fish command");
 *   cli.addNode(fishList);
 *   CliNode fishMarlinCommand = CliNode.makeCommand("marlin", "This is the marlin command", new CliAction() { ... });
 *   fishList.addNode(fishMarlinCommand);
 *   ...
 *   cli.parse(args);
 * </pre>
 *
 * The above would result in the following command line API.
 * <pre>
 *   root-command [-f] [-v|--verbose] [-o|--output location] tuna [-a]
 *   root-command [-f] [-v|--verbose] [-o|--output location] fish marlin
 * </pre>
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class Cli {

  private CliNode root;

  public Cli(String name, String description) {
    root = CliNode.makeList(name, description);
  }

  public void addFlag(CliFlag flag) {
    root.addFlag(flag);
  }

  public void addNode(CliNode node) {
    root.addNode(node);
  }

  public void parse(String[] args) {
    root.parse(args);
  }

  public String getUsage() {
    return root.getUsage();
  }
}
