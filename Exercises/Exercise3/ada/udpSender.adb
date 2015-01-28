with Ada.Text_IO; use Ada.Text_IO;
with Ada.Strings.Unbounded; use Ada.Strings.Unbounded;
with GNAT.Sockets; use GNAT.Sockets;
with Ada.Streams;
use type Ada.Streams.Stream_Element_Count;

procedure udpsender is
  client  : Socket_Type; address : Sock_Addr_Type;
  Channel : Stream_Access;

  Message : Unbounded_String;

  Offset : Ada.Streams.Stream_Element_Count;
  Data   : Ada.Streams.Stream_Element_Array (1 .. 256);
  Last	 : Ada.Streams.Stream_Element_Offset;
  sendData: Ada.Streams.Stream_Element_Array := (1, 2, 3);

begin

  Initialize(Process_Blocking_IO => False);
  --GNAT.Sockets.Initialize;  -- initialize a specific package
  Create_Socket(client,  Family_Inet, Socket_Datagram);

  address.Port := 20008;

  address.Addr := Inet_Addr("129.241.187.136");
  Bind_Socket (client, address);

  Channel := Stream (client, address);
  Ada.Text_IO.Put_Line ("Created Stream");

   loop
      Last := sendData'last;
      Send_Socket(client, sendData, Last,To => address);
      Ada.Streams.Read (Channel.All, Data, Offset);
      address := Get_Address (Channel);
      Ada.Text_IO.Put_Line ("from " & Image (address) & " : ");
      exit when Offset = 0;
      for I in 1 .. Offset loop
         Ada.Text_IO.Put (Character'Val (Data(I)));
      end loop;
      Put_Line("");
   end loop;

  Close_Socket(client);

end udpsender;
