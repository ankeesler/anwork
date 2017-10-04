package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.app.cli.CliActionCreator;
import com.marshmallow.anwork.app.cli.CliCommand;

import java.util.HashMap;
import java.util.Map;

/**
 * This is a {@link CliActionCreator} that creates {@link TestCliAction} instances and stores them
 * in a <code>static</code> {@link Map}.
 *
 * <p>
 * Note that the methods offered on this class are <code>static</code> and they track <em>all</em>
 * of the {@link TestCliAction}'s that this {@link CliActionCreator} has created across
 * <em>all</em> instances of this class.
 * </p>
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TestCliActionCreator implements CliActionCreator {

  private static Map<String, TestCliAction> createdActions = new HashMap<String, TestCliAction>();

  /**
   * Forget about all of the {@link CliAction}'s that this class has created.
   */
  public static void resetCreatedActions() {
    createdActions.clear();
  }

  /**
   * Get the {@link TestCliAction} that was created for the provided {@link CliCommand} name.
   *
   * @param commandName The name of the {@link CliCommand}, e.g., "marlin", "show",
   *     "bring-home-bacon", etc.
   * @return The {@link TestCliAction} that was created for the provided {@link CliCommand} name,
   *     or <code>null</code> if there is no known {@link TestCliAction} for the provided
   *     {@link CliCommand} name.
   */
  public static TestCliAction getCreatedAction(String commandName) {
    return createdActions.get(commandName);
  }

  @Override
  public CliAction createAction(String commandName) {
    TestCliAction action = new TestCliAction();
    createdActions.put(commandName, action);
    return action;
  }
}
