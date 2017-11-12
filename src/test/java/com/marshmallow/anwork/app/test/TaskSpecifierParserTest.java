package com.marshmallow.anwork.app.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.app.TaskSpecifierParser;
import com.marshmallow.anwork.task.TaskManager;

import org.junit.Before;
import org.junit.Test;

/**
 * This is a unit test for the {@link TaskSpecifierParser}.
 *
 * <p>
 * Created Nov 12, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskSpecifierParserTest {

  private static final String TASK_A_NAME = "task-a";
  private static final String TASK_B_NAME = "task-b";

  private static final String NONEXISTENT_TASK_NAME = "task-nonexistent";
  private static final int NONEXISTENT_TASK_ID = 1089723;

  private TaskManager manager;
  private TaskSpecifierParser parser;
  private int taskAId;
  private int taskBId;

  /**
   * This method initializes the {@link TaskManager} with some {@link Task}'s and initializes the
   * {@link ArgumentType} under test to use that {@link TaskManager}.
   */
  @Before
  public void setupTaskManager() {
    manager = new TaskManager();
    manager.createTask(TASK_A_NAME, "", 1);
    manager.createTask(TASK_B_NAME, "", 2);

    // Make sure the expected task IDs are correct.
    taskAId = manager.getTasks()[0].getId();
    taskBId = manager.getTasks()[1].getId();

    parser = new TaskSpecifierParser(manager);
  }

  @Test
  public void testGoodTaskNameConversion() {
    String[] converted = parser.parse(TASK_A_NAME);
    assertArrayEquals(new String[] { TASK_A_NAME }, converted);
    converted = parser.parse(TASK_B_NAME);
    assertArrayEquals(new String[] { TASK_B_NAME }, converted);
  }

  @Test
  public void testGoodTaskIdConversion() {
    String taskAIdString = String.format("%c%d", TaskSpecifierParser.SPECIAL_CHARACTER, taskAId);
    String[] converted = parser.parse(taskAIdString);
    assertArrayEquals(new String[] { TASK_A_NAME }, converted);
    String taskBIdString = String.format("%c%d", TaskSpecifierParser.SPECIAL_CHARACTER, taskBId);
    converted = parser.parse(taskBIdString);
    assertArrayEquals(new String[] { TASK_B_NAME }, converted);
  }

  @Test(expected = IllegalArgumentException.class)
  public void test0LengthString() {
    parser.parse("");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNonexistentTask() {
    assertNull(parser.parse(NONEXISTENT_TASK_NAME));
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNonexistentId() {
    String badSpecialSignFormat = String.format("%c%d",
                                                TaskSpecifierParser.SPECIAL_CHARACTER,
                                                NONEXISTENT_TASK_ID);
    parser.parse(badSpecialSignFormat);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadSpecialSignFormat() {
    String badSpecialSignFormat = String.format("%c%c%d",
                                                TaskSpecifierParser.SPECIAL_CHARACTER,
                                                TaskSpecifierParser.SPECIAL_CHARACTER,
                                                taskAId);
    parser.parse(badSpecialSignFormat);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadSpecialSignFormatWithTaskName() {
    String badSpecialSignFormat = String.format("%c%s",
                                                TaskSpecifierParser.SPECIAL_CHARACTER,
                                                TASK_A_NAME);
    parser.parse(badSpecialSignFormat);
  }
}
