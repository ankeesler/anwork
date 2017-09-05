package com.marshmallow.anwork.task.test;

import org.junit.Test;

import com.marshmallow.anwork.core.test.SerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskState;

/**
 * A {@link SerializerTest<T>} for {@link Task} objects.
 *
 * @author Andrew
 * @date Sep 4, 2017
 */
public class TaskSerializerTest extends SerializerTest<Task> {

  public TaskSerializerTest() {
    super(Task.serializer());
  }

  @Test
  public void emptyStringTest() {
    assertBad("");
  }

  @Test
  public void testJustPlainWrongTask() {
    assertBad("This is not a task");
  }

  @Test
  public void badPrefixTest() {
    assertBad("name=a;id=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Tas:name=a;id=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Taskname=a;id=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Task;name=a;id=0;description=b;date=123;priority=1;state=WAITING;");
  }

  @Test
  public void missingFieldsTest() {
    assertBad("Task:id=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;date=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b;date=123;priority=1;");
  }

  @Test
  public void missingDeliminatorsTest() {
    assertBad("Task:name=aid=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;description=bdate=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b;date=123;priority=1;state=WAITING");
  }

  @Test
  public void tooManyDeliminatorsTest() {
    assertBad("Task:name=a;;id=0;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b;date=123;;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b;date=123;priority=1;state=WAITING;;");
  }

  @Test
  public void spacesTest() {
    assertBad("Task :name=a;id=0;description=b;date=123;priority=3;state=WAITING;");
    assertBad("Task: name=a;id=0;description=b;date=123;priority=3;state=WAITING;");
    assertBad("Task:name=a; id=0;description=b;date=123;priority=3;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b; date=123;priority=3;state=WAITING;");
    assertBad("Task:name=a;id=0;description=b; date=123;priority=3;state=WAITING; ");
  }

  @Test
  public void outOfOrderTest() {
    assertBad("Task:name=a;description=b;id=0;date=123;priority=3;state=WAITING;");
    assertBad("Task:id=1;name=a;description=b;date=foo;priority=1;state=WAITING;");
    assertBad("Task:id=1;name=a;description=b;date=foo;state=WAITING;priority=1;");
  }

  @Test
  public void badNumbersTest() {
    assertBad("Task:name=a;id=whatever;description=b;date=123;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=1;description=b;date=foo;priority=1;state=WAITING;");
    assertBad("Task:name=a;id=1;description=b;date=123;priority=blah;state=WAITING;");
  }

  @Test
  public void datesTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=%d;priority=3;state=WAITING;";
    for (long i = 1; i < Integer.MAX_VALUE; i <<= 1) {
      assertGood(String.format(goodTaskFormat, i));
    }
  }

  @Test
  public void priorityTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=123;priority=%d;state=WAITING;";
    for (long i = 1; i < Integer.MAX_VALUE; i *= 3) {
      assertGood(String.format(goodTaskFormat, i));
    }
  }

  @Test
  public void statesTest() {
    String goodTaskFormat = "Task:name=a;id=0;description=b;date=123;priority=3;state=%s;";
    for (TaskState taskState : TaskState.values()) {
      assertGood(String.format(goodTaskFormat, taskState.name()));
    }
  }
}