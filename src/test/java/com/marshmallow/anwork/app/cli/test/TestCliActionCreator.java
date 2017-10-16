package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ActionCreator;
import com.marshmallow.anwork.app.cli.Command;

import java.util.HashMap;
import java.util.Map;

/**
 * This is a {@link ActionCreator} that creates {@link TestCliAction} instances and stores them
 * in a <code>static</code> {@link Map}.
 *
 * <p>
 * Note that the methods offered on this class are <code>static</code> and they track <em>all</em>
 * of the {@link TestCliAction}'s that this {@link ActionCreator} has created across
 * <em>all</em> instances of this class.
 * </p>
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TestCliActionCreator implements ActionCreator {

  private static Map<String, TestCliAction> createdActions = new HashMap<String, TestCliAction>();

  /**
   * Forget about all of the {@link Action}'s that this class has created.
   */
  public static void resetCreatedActions() {
    createdActions.clear();
  }

  /**
   * Get the {@link TestCliAction} that was created for the provided {@link Command} name.
   *
   * @param commandName The name of the {@link Command}, e.g., "marlin", "show",
   *     "bring-home-bacon", etc.
   * @return The {@link TestCliAction} that was created for the provided {@link Command} name,
   *     or <code>null</code> if there is no known {@link TestCliAction} for the provided
   *     {@link Command} name.
   */
  public static TestCliAction getCreatedAction(String commandName) {
    return createdActions.get(commandName);
  }

  @Override
  public Action createAction(String commandName) {
    TestCliAction action = new TestCliAction();
    createdActions.put(commandName, action);
    return action;
  }
}
