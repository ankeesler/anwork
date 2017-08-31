package com.marshmallow.anwork.core.test;

import org.junit.Test;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskState;

import static org.junit.Assert.*;

/**
 * This is a generic test for a {@link Serializer<Task>}.
 *
 * @author Andrew
 * @date Aug 31, 2017
 */
public class SerializerTest {

  @Test
  public void emptyStringTest() {
    runSerializationTest("", false);
  }

  @Test
  public void testJustPlainWrongTask() {
    runSerializationTest("This is not a task", false);
  }

  @Test
  public void badPrefixTest() {
    runSerializationTest("name=a;id=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Tas:name=a;id=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Taskname=a;id=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task;name=a;id=0;description=b;date=123;priority=1;state=WAITING;", false);
  }

  @Test
  public void missingFieldsTest() {
    runSerializationTest("Task:id=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;description=b;date=123;priority=1;", false);
  }

  @Test
  public void missingDeliminatorsTest() {
    runSerializationTest("Task:name=aid=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;description=bdate=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;description=b;date=123;priority=1;state=WAITING", false);
  }

  @Test
  public void tooManyDeliminatorsTest() {
    runSerializationTest("Task:name=a;;id=0;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;description=b;date=123;;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=0;description=b;date=123;priority=1;state=WAITING;;", false);
  }

  @Test
  public void outOfOrderTest() {
    runSerializationTest("Task:name=a;description=b;id=0;date=123;priority=3;state=WAITING;", false);
    runSerializationTest("Task:id=1;name=a;description=b;date=foo;priority=1;state=WAITING;", false);
    runSerializationTest("Task:id=1;name=a;description=b;date=foo;state=WAITING;priority=1;", false);
  }

  @Test
  public void badNumbersTest() {
    runSerializationTest("Task:name=a;id=whatever;description=b;date=123;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=1;description=b;date=foo;priority=1;state=WAITING;", false);
    runSerializationTest("Task:name=a;id=1;description=b;date=123;priority=blah;state=WAITING;", false);
  }

  @Test
  public void datesTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=%d;priority=3;state=WAITING;";
    for (long i = 1; i < Integer.MAX_VALUE; i <<= 1) {
      runSerializationTest(String.format(goodTaskFormat, i), true);
    }
  }

  @Test
  public void priorityTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=123;priority=%d;state=WAITING;";
    for (long i = 1; i < Integer.MAX_VALUE; i *= 3) {
      runSerializationTest(String.format(goodTaskFormat, i), true);
    }
  }

  @Test
  public void statesTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=123;priority=3;state=%s;";
    for (TaskState taskState : TaskState.values()) {
      runSerializationTest(String.format(goodTaskFormat, taskState.name()), true);
    }
  }

  private static void runSerializationTest(String string, boolean success) {
    Serializer<Task> taskSerializer = Task.serializer();

    Task task = taskSerializer.unmarshall(string);
    String message = String.format("Expected %s task from %s", (success ? "a" : "no"), string);
    assertEquals(message, task != null, success);

    if (success) {
      message = String.format("Marshalled task (%s) does not match original string (%s)", task, string);
      assertEquals(message, string, taskSerializer.marshall(task));
    }
  }
}
