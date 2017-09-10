package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.SerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

import org.junit.BeforeClass;
import org.junit.Test;

/**
 * A {@link SerializerTest} for {@link TaskManager} objects.
 *
 * @author Andrew
 * Created Sep 4, 2017
 */
public class TaskManagerSerializerTest extends SerializerTest<TaskManager> {

  private static final String GOOD_TASK_A
    = "Task:name=a;id=0;description=b;date=123;priority=3;state=WAITING;";
  private static final String GOOD_TASK_B
    = "Task:name=b;id=0;description=b;date=123;priority=3;state=WAITING;";
  private static final String GOOD_TASK_C
    = "Task:name=c;id=0;description=b;date=123;priority=3;state=WAITING;";
  private static final String BAD_TASK_A
    = "Task:name=a;id=0;escription=b;date=123;priority=3;state=WAITING;";

  private static Task taskA;
  private static Task taskB;
  private static Task taskC;

  public TaskManagerSerializerTest() {
    super(TaskManager.serializer());
  }

  @BeforeClass
  public static void setupTasks() {
    Serializer<Task> taskSerializer = Task.serializer();
    taskA = taskSerializer.unmarshall(GOOD_TASK_A);
    assertNotNull(taskA);
    taskB = taskSerializer.unmarshall(GOOD_TASK_B);
    assertNotNull(taskB);
    taskC = taskSerializer.unmarshall(GOOD_TASK_C);
    assertNotNull(taskC);
  }

  @Test
  public void testEmpty() {
    assertBad("");
  }

  @Test
  public void testJustPlainWrongTask() {
    assertBad("This is not a task");
  }

  @Test
  public void testBadStart() {
    assertBad(":" + GOOD_TASK_A);
    assertBad("Task:" + GOOD_TASK_A);
    assertBad("TaskManage:" + GOOD_TASK_A);
    assertBad("TaskManager;" + GOOD_TASK_A);
  }

  @Test
  public void testBadTask() {
    assertBad("TaskManager:" + BAD_TASK_A);
    assertBad("TaskManager:" + GOOD_TASK_A + "," + BAD_TASK_A);
  }

  @Test
  public void testBadSeparator() {
    assertBad("TaskManager:" + GOOD_TASK_A);
    assertBad("TaskManager:" + GOOD_TASK_A + GOOD_TASK_B);
    assertBad("TaskManager:" + GOOD_TASK_A + GOOD_TASK_B + ",");
  }

  @Test
  public void testBadSpaces() {
    assertBad("TaskManager :" + GOOD_TASK_A + ",");
    assertBad("TaskManager: " + GOOD_TASK_A + ",");
    assertBad("TaskManager:" + GOOD_TASK_A + " ,");
  }

  @Test
  public void testNoTasks() {
    TaskManager manager = assertGood("TaskManager:");
    assertNull(manager.getCurrentTask());
    assertEquals(0, manager.getTaskCount());
  }

  @Test
  public void testSingleTask() {
    TaskManager manager = assertGood("TaskManager:" + GOOD_TASK_A + ",");
    assertNull(manager.getCurrentTask());
    assertEquals(taskA.getState().name().toLowerCase(), manager.getState(taskA.getName()));
    assertEquals(1, manager.getTaskCount());
  }

  @Test
  public void testMultipleTasks() {
    TaskManager manager
      = assertGood("TaskManager:" + GOOD_TASK_A + ",*" + GOOD_TASK_B + "," + GOOD_TASK_C + ",");
    assertEquals(taskB.getName(), manager.getCurrentTask());
    assertEquals(taskA.getState().name().toLowerCase(), manager.getState(taskA.getName()));
    assertEquals(taskB.getState().name().toLowerCase(), manager.getState(taskB.getName()));
    assertEquals(taskC.getState().name().toLowerCase(), manager.getState(taskC.getName()));
    assertEquals(3, manager.getTaskCount());
  }
}
