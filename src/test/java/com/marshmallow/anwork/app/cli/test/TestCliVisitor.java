package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.CliVisitor;

import java.util.ArrayList;
import java.util.List;

/**
 * This is a {@link CliVisitor} that is used to test the CLI visitor functionality.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class TestCliVisitor implements CliVisitor {

  private List<String> visitedShortFlags = new ArrayList<String>();
  private List<String> visitedShortFlagsWithParameters = new ArrayList<String>();
  private List<String> visitedLongFlags = new ArrayList<String>();
  private List<String> visitedLongFlagsWithParameters = new ArrayList<String>();
  private List<String> visitedCommands = new ArrayList<String>();
  private List<String> visitedLists = new ArrayList<String>();

  /**
   * Get the short flags that were visited by this {@link CliVisitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedShortFlags() {
    return visitedShortFlags;
  }

  /**
   * Get the short flags with parameters that were visited by this {@link CliVisitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedShortFlagsWithParameters() {
    return visitedShortFlagsWithParameters;
  }

  /**
   * Get the long flags that were visited by this {@link CliVisitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedLongFlags() {
    return visitedLongFlags;
  }

  /**
   * Get the long flags with parameters that were visited by this {@link CliVisitor}.
   *
   * <p>
   * Note that the order of the flags may not match the order in which they were added!
   * </p>
   *
   * @return The flags that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedLongFlagsWithParameters() {
    return visitedLongFlagsWithParameters;
  }

  /**
   * Get the commands that were visited by this {@link CliVisitor}. Note that the order of the
   * commands may not match the order in which they were added to a list!
   *
   * @return The commands that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedCommands() {
    return visitedCommands;
  }

  /**
   * Get the lists that were visited by this {@link CliVisitor}. Note that the order of the lists
   * may not match the order in which they were added to a list!
   *
   * @return The lists that were visited by this {@link CliVisitor}
   */
  public List<String> getVisitedLists() {
    return visitedLists;
  }

  @Override
  public void visitShortFlag(String shortFlag,
                             String description) {
    visitedShortFlags.add(shortFlag);
  }

  @Override
  public void visitShortFlagWithParameter(String shortFlag,
                                          String parameterName,
                                          String description) {
    visitedShortFlagsWithParameters.add(shortFlag);
  }

  @Override
  public void visitLongFlag(String shortFlag,
                            String longFlag,
                            String description) {
    visitedLongFlags.add(shortFlag);
  }

  @Override
  public void visitLongFlagWithParameter(String shortFlag,
                                         String longFlag,
                                         String parameterName,
                                         String description) {
    visitedLongFlagsWithParameters.add(shortFlag);
  }

  @Override
  public void visitList(String name,
                        String description) {
    visitedLists.add(name);
  }

  @Override
  public void visitCommand(String name,
                           String description) {
    visitedCommands.add(name);
  }

}
