package com.marshmallow.anwork.core.test;

import org.junit.Test;

import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Serializer;

import static org.junit.Assert.*;

import java.io.File;
import java.io.IOException;
import java.util.Arrays;
import java.util.Collection;
import java.util.Collections;

/**
 * This is a unit test for a {@link FilePersister}.
 *
 * @author Andrew
 * @date Sep 6, 2017
 */
public class FilePersisterTest {

  private static final File TEST_RESOURCE_ROOT = new File(TestUtilities.TEST_RESOURCES_ROOT, "file-persister-test");

  private static final String DEFAULT_CONTEXT = "default-context";
  private static final String CONTEXT_A = "context-a";
  private static final String CONTEXT_B = "context-b";
  private static final Serializer<Student> DEFAULT_SERIALIZER = Student.serializer();
  private static final FilePersister<Student> DEFAULT_PERSISTER = new FilePersister<Student>(TEST_RESOURCE_ROOT);

  @Test(expected = IOException.class)
  public void testContextDoesNotExist() throws IOException {
    String context = "this-context-does-not-exist";
    assertFalse(DEFAULT_PERSISTER.contextExists(context));
    DEFAULT_PERSISTER.load(context, DEFAULT_SERIALIZER);
  }

  @Test(expected = IOException.class)
  public void testBadSerialization() throws IOException {
    String context = "this-context-has-a-bad-serialization";
    assertTrue(DEFAULT_PERSISTER.contextExists(context));
    DEFAULT_PERSISTER.load(context, DEFAULT_SERIALIZER);
  }

  @Test()
  public void loadOneTest() throws IOException {
    String context = "load-one-test";
    assertTrue(DEFAULT_PERSISTER.contextExists(context));
    Collection<Student> loadeds = DEFAULT_PERSISTER.load(context, DEFAULT_SERIALIZER);
    assertEquals(1, loadeds.size());

    Student loaded = loadeds.toArray(new Student[0])[0];
    assertEquals("Tuna Fish Marlin", loaded.getName());
    assertEquals(8675309, loaded.getId());
  }

  @Test()
  public void loadThreeTest() throws IOException {
    String context = "load-three-test";
    assertTrue(DEFAULT_PERSISTER.contextExists(context));
    Collection<Student> loadeds = DEFAULT_PERSISTER.load(context, DEFAULT_SERIALIZER);
    assertEquals(3, loadeds.size());

    Student[] loadedsArray = loadeds.toArray(new Student[0]);
    assertEquals("Tuna", loadedsArray[0].getName());
    assertEquals(100, loadedsArray[0].getId());
    assertEquals("Fish", loadedsArray[1].getName());
    assertEquals(200, loadedsArray[1].getId());
    assertEquals("Marlin", loadedsArray[2].getName());
    assertEquals(300, loadedsArray[2].getId());
  }

  @Test
  public void saveOneTest() throws IOException {
    assertTrue(DEFAULT_PERSISTER.contextExists(DEFAULT_CONTEXT));
    Student saved = new Student("Andrew", 12345);
    DEFAULT_PERSISTER.save(DEFAULT_CONTEXT, DEFAULT_SERIALIZER, Collections.singleton(saved));

    Collection<Student> loadeds = DEFAULT_PERSISTER.load(DEFAULT_CONTEXT, DEFAULT_SERIALIZER);
    assertEquals(1, loadeds.size());

    Student loaded = loadeds.toArray(new Student[0])[0];
    assertEquals("Andrew", loaded.getName());
    assertEquals(12345, loaded.getId());
  }

  @Test
  public void saveThreeTest() throws IOException {
    assertTrue(DEFAULT_PERSISTER.contextExists(DEFAULT_CONTEXT));
    Collection<Student> saveds = Arrays.asList(new Student("Andrew", 1), new Student("AC", 2), new Student("Mom", 3));
    DEFAULT_PERSISTER.save(DEFAULT_CONTEXT, DEFAULT_SERIALIZER, saveds);

    Collection<Student> loadeds = DEFAULT_PERSISTER.load(DEFAULT_CONTEXT, DEFAULT_SERIALIZER);
    assertEquals(3, loadeds.size());

    Student[] loadedsArray = loadeds.toArray(new Student[0]);
    assertEquals("Andrew", loadedsArray[0].getName());
    assertEquals(1, loadedsArray[0].getId());
    assertEquals("AC", loadedsArray[1].getName());
    assertEquals(2, loadedsArray[1].getId());
    assertEquals("Mom", loadedsArray[2].getName());
    assertEquals(3, loadedsArray[2].getId());
  }

  @Test
  public void doubleContextTest() throws IOException {
    assertTrue(DEFAULT_PERSISTER.contextExists(CONTEXT_A));
    assertTrue(DEFAULT_PERSISTER.contextExists(CONTEXT_B));
    Collection<Student> savedA = Arrays.asList(new Student("Andrew", 1));
    Collection<Student> savedB = Arrays.asList(new Student("AC", 2), new Student("Mom", 3));
    DEFAULT_PERSISTER.save(CONTEXT_A, DEFAULT_SERIALIZER, savedA);
    DEFAULT_PERSISTER.save(CONTEXT_B, DEFAULT_SERIALIZER, savedB);

    Collection<Student> loadedsA = DEFAULT_PERSISTER.load(CONTEXT_A, DEFAULT_SERIALIZER);
    Collection<Student> loadedsB = DEFAULT_PERSISTER.load(CONTEXT_B, DEFAULT_SERIALIZER);
    assertEquals(1, loadedsA.size());
    assertEquals(2, loadedsB.size());

    Student[] loadedsArrayA = loadedsA.toArray(new Student[0]);
    Student[] loadedsArrayB = loadedsB.toArray(new Student[0]);
    assertEquals("Andrew", loadedsArrayA[0].getName());
    assertEquals(1, loadedsArrayA[0].getId());
    assertEquals("AC", loadedsArrayB[0].getName());
    assertEquals(2, loadedsArrayB[0].getId());
    assertEquals("Mom", loadedsArrayB[1].getName());
    assertEquals(3, loadedsArrayB[1].getId());
  }
}