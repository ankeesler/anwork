package com.marshmallow.anwork.app;

import java.io.FileNotFoundException;
import java.io.PrintWriter;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliVisitor;

/**
 * This is an app that generates the documentation for the ANWORK CLI from the
 * {@link AnworkCliCreator}.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkCliDocumentationGenerator implements CliVisitor {

  private static final String FILENAME = "CLI.md";

  private static enum State {
    NONE,
    FLAG,
    LIST,
    COMMAND,
    ;
  }

  /**
   * This is the main method for the {@link AnworkCliDocumentationGenerator} app.
   *
   * @param args The command line arguments for this app.
   */
  public static void main(String[] args) {
    try (PrintWriter writer = new PrintWriter(FILENAME)) {
      AnworkAppConfig config = new AnworkAppConfig();
      Cli cli = new AnworkCliCreator(config).makeCli();
      CliVisitor visitor = new AnworkCliDocumentationGenerator(writer);
      cli.visit(visitor);
    } catch (FileNotFoundException e) {
      System.out.println("Error: " + e);
    }
  }

  private PrintWriter writer;
  private State state;

  private AnworkCliDocumentationGenerator(PrintWriter writer) {
    this.writer = writer;
    this.state = State.NONE;

    writer.println("This documentation is generated from " + this.getClass().getName());
    writer.println();
  }

  private void checkFlagState() {
    if (state != State.FLAG) {
      writer.println("## Flags");
      state = State.FLAG;
    }
  }

  private void checkListState() {
    if (state != State.LIST) {
      state = State.LIST;
    }
  }

  private void checkCommandState() {
    if (state != State.COMMAND) {
      writer.println("## Commands");
      state = State.COMMAND;
    }
  }

  @Override
  public void visitShortFlag(String shortFlag, String description) {
    checkFlagState();
    writer.println(String.format("- %s: %s", shortFlag, description));
  }

  @Override
  public void visitShortFlagWithParameter(String shortFlag,
                                          String parameterName,
                                          String description) {
    checkFlagState();
    writer.println(String.format("- -%s <%s>: %s", shortFlag, parameterName, description));
  }

  @Override
  public void visitLongFlag(String shortFlag,
                            String longFlag,
                            String description) {
    checkFlagState();
    writer.println(String.format("- -%s|--%s: %s", shortFlag, longFlag, description));
  }

  @Override
  public void visitLongFlagWithParameter(String shortFlag,
                                         String longFlag,
                                         String parameterName,
                                         String description) {
    checkFlagState();
    writer.println(String.format("- -%s|--%s <%s>: %s",
                                 shortFlag,
                                 longFlag,
                                 parameterName,
                                 description));
  }

  @Override
  public void visitList(String name, String description) {
    checkListState();
    writer.println(String.format("# %s: %s", name, description));
  }

  @Override
  public void visitCommand(String name, String description) {
    checkCommandState();
    writer.println(String.format("- %s: %s", name, description));
  }
}
